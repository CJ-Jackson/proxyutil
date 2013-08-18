package proxyutil

import (
	//"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Set up Reverse Proxy
func ReverseProxy(rawurl string) http.Handler {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

/*
func ReverseProxyUnix(filepath string) http.Handler {
	proxy := &httputil.ReverseProxy{}

	proxy.Director = func(r *http.Request) {
		target, _ := url.Parse("http://" + r.Host)

		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
	}

	dial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", filepath)
	}

	proxy.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, Dial: dial}

	return proxy
}
*/
