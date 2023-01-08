package irankish

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

// TODO: Import httpjson package to handle I/O using json

// Proxy is to run a http proxy server to redirect requests
//	This is because there's IP limit on payment gateways
//	Each payment gateway is only allowed to be accessed by specific ip addresses
//	You can run this proxy server in an allowed IP address and initiate your gateway using WithProxy
type Proxy struct {
	httpUri string
	mux     *http.ServeMux

	homePageHandler func(writer http.ResponseWriter, r *http.Request)

	targetUrl          *url.URL
	callbackUrlsLocker sync.Mutex
	callbackUrls       map[string]*url.URL

	logger Logger
}

func NewProxy(httpUri string, logger Logger) *Proxy {
	return &Proxy{httpUri: httpUri, mux: http.NewServeMux(), targetUrl: host, logger: logger}
}

func NewProxyWithMux(httpUri string, mux *http.ServeMux, logger Logger) *Proxy {
	return &Proxy{httpUri: httpUri, mux: mux, targetUrl: host, logger: logger}
}

// EnableCallbackUrls by calling this method /add_callback_url endpoint will be activated
//	This is to proxy back callback (revert url) requests to specified endpoints
func (p *Proxy) EnableCallbackUrls() *Proxy {
	p.callbackUrlsLocker.Lock()
	defer p.callbackUrlsLocker.Unlock()

	p.callbackUrls = map[string]*url.URL{}
	p.mux.HandleFunc("/add_callback_url", p.handleAddCallback)

	return p
}

// SetTargetUrl is used to replace payment gateway's base url
func (p *Proxy) SetTargetUrl(target *url.URL) *Proxy {
	p.targetUrl = target

	return p
}

func (p *Proxy) SetHomePageHandler(handler func(writer http.ResponseWriter, r *http.Request)) *Proxy {
	p.homePageHandler = handler

	return p
}

func (p *Proxy) Start() error {
	p.registerRoutes()

	return http.ListenAndServe(p.httpUri, p.mux)
}

func (p *Proxy) registerRoutes() {
	p.mux.HandleFunc("/redirect", p.handleRedirect)
	p.mux.HandleFunc("/", p.handleRequests)
}

func (p *Proxy) handleRequests(writer http.ResponseWriter, request *http.Request) {
	if request.URL.String() == "/" && p.homePageHandler != nil {
		p.homePageHandler(writer, request)
		return
	}

	defer request.Body.Close()

	b, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	rr, err := http.NewRequest(request.Method, request.URL.String(), bytes.NewReader(b))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	rr.Header.Set("Content-Type", "application/json")

	// Fun fact is that if there are any other headers like x-forwarded-server is set, it will complain about the ip
	// Maybe it's possible to bypass ip by setting x-forwarded-for header.
	// rr.Header = request.Header <--- Removed for above mentioned reason

	// If user has assigned a callback for this endpoint, forward the request there, else forward it to payment gateway
	//	This is to proxy requests coming from payment gateway to custom urls
	forwardTo := p.targetUrl
	if callback := p.getEndpointCallback(request.URL.String()); callback != nil {
		forwardTo = callback
	}

	if p.logger != nil {
		go p.logger.Println(fmt.Sprintf("Forwarding from %s %s to %s", request.RemoteAddr, request.URL.String(), forwardTo.String()))
	}

	httputil.NewSingleHostReverseProxy(forwardTo).ServeHTTP(writer, rr)
}

func (p *Proxy) getEndpointCallback(endpoint string) *url.URL {
	p.callbackUrlsLocker.Lock()
	defer p.callbackUrlsLocker.Unlock()

	if callback, ok := p.callbackUrls[endpoint]; ok {
		return callback
	}

	return nil
}

func (p *Proxy) handleRedirect(writer http.ResponseWriter, request *http.Request) {
	token := request.URL.Query().Get("token")
	if token == "" {
		_, _ = writer.Write([]byte("invalid_request"))
		return
	}

	mkr := &MakeTokenResult{Token: token}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(mkr.RedirectForm()))
}

// handleAddCallback receives the request as POST form with following fields
//	- endpoint
//	- callback_url
//	Then registers endpoint with the mux and forwards incoming callbacks from payment gateway to callback_url
func (p *Proxy) handleAddCallback(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	endpoint := request.FormValue("endpoint")
	callbackUrl := request.FormValue("callback_url")

	if endpoint == "" || callbackUrl == "" {
		_, _ = writer.Write([]byte("empty_values"))
		return
	}

	callbackUrlParsed, err := url.Parse(callbackUrl)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	p.callbackUrlsLocker.Lock()
	defer p.callbackUrlsLocker.Unlock()

	// TODO: Add a mechanism to show error if the user tries to register a duplicate endpoint that they don't own
	p.callbackUrls[endpoint] = callbackUrlParsed
}

func AddCallbackUrl(serverAddress, endpoint, callbackUrl string) error {
	values := url.Values{}
	values.Set("endpoint", endpoint)
	values.Set("callback_url", callbackUrl)

	req, err := http.NewRequest("POST", serverAddress+"/add_callback_url", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
