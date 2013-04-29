package proxyutil

import (
	"net/http/httputil"
	"net/url"
)

func ReverseProxy(rawurl string) *httputil.ReverseProxy {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}
