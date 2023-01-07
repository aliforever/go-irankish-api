package irankish

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
)

// Proxy is to run a http proxy server to redirect requests
//	This is because there's IP limit on payment gateways
//	Each payment gateway is only allowed to be accessed by specific ip addresses
//	You can run this proxy server in an allowed IP address and initiate your gateway using WithProxy
type Proxy struct {
	httpUri string
	mux     *http.ServeMux
}

func NewProxy(httpUri string) *Proxy {
	return &Proxy{httpUri: httpUri, mux: http.NewServeMux()}
}

func NewProxyWithMux(httpUri string, mux *http.ServeMux) *Proxy {
	return &Proxy{httpUri: httpUri, mux: mux}
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

	httputil.NewSingleHostReverseProxy(host).ServeHTTP(writer, rr)
}
