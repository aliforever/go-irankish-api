package irankish

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	callbackUrlsLocker sync.Mutex
	callbackUrls       map[string]*url.URL
}

func NewProxy(httpUri string) *Proxy {
	return &Proxy{httpUri: httpUri, mux: http.NewServeMux()}
}

func NewProxyWithMux(httpUri string, mux *http.ServeMux) *Proxy {
	return &Proxy{httpUri: httpUri, mux: mux}
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

func (p *Proxy) Start() error {
	p.registerRoutes()

	return http.ListenAndServe(p.httpUri, p.mux)
}

func (p *Proxy) registerRoutes() {
	p.mux.HandleFunc("/", p.handleRequests)
}

func (p *Proxy) handleRequests(writer http.ResponseWriter, request *http.Request) {
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

	rr.Header = request.Header

	// If user has assigned a callback for this endpoint, forward the request there, else forward it to payment gateway
	//	This is to proxy requests coming from payment gateway to custom urls
	forwardTo := host
	if callback := p.getEndpointCallback(request.URL.String()); callback != nil {
		forwardTo = callback
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