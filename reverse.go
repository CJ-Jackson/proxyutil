package proxyutil

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ReverseProxy(rawurl string) http.Handler {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}
