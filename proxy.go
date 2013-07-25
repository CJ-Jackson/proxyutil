package proxyutil

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Domain interface {
	MatchString(string) bool
}

type strDomain string

func (s strDomain) MatchString(str string) bool {
	return string(s) == str
}

type strDomains []string

func (s strDomains) MatchString(str string) bool {
	for _, ss := range []string(s) {
		if ss == str {
			return true
		}
	}
	return false
}

func Str(str ...string) Domain {
	if len(str) <= 0 {
		return nil
	}
	if len(str) == 1 {
		return strDomain(str[0])
	}
	return strDomains(str)
}

// Alais of Str
func Names(names ...string) Domain {
	return Str(names...)
}

func RegExp(pattern string) Domain {
	return regexp.MustCompile(pattern)
}

type Host struct {
	Domain  Domain
	Prepare http.Handler
	Proxy   http.Handler
	WS      http.Handler
	Finish  http.Handler
}

type ListOfHosts struct {
	sync.Mutex
	li []Host
}

func (li *ListOfHosts) Register(h ...Host) {
	li.Lock()
	defer li.Unlock()
	li.li = append(li.li, h...)
}

func (li *ListOfHosts) getHost() []Host {
	li.Lock()
	defer li.Unlock()
	return append([]Host{}, li.li...)
}

func hostnameWithoutPort(str string) string {
	hostname, _, err := net.SplitHostPort(str)
	if err != nil {
		hostname = strings.Split(str, ":")[0]
	}
	return hostname
}

func hostDealer(host *Host, res http.ResponseWriter, req *http.Request) {
	if IsWebSocket(req) {
		if host.WS != nil {
			host.WS.ServeHTTP(res, req)
		}
		return
	}

	if host.Prepare != nil {
		host.Prepare.ServeHTTP(res, req)
	}

	if host.Proxy != nil {
		host.Proxy.ServeHTTP(res, req)
	}

	if host.Finish != nil {
		host.Finish.ServeHTTP(res, req)
	}
}

func proxy(res http.ResponseWriter, req *http.Request) {
	hostname := hostnameWithoutPort(req.Host)
	if Prepare != nil {
		Prepare.ServeHTTP(res, req)
	}

	if Finish != nil {
		defer Finish.ServeHTTP(res, req)
	}

	cache.RLock()
	switch host := cache.m[hostname].(type) {
	case *Host:
		cache.Unlock()
		hostDealer(host, res, req)
		return
	case bool:
		cache.RUnlock()
		return
	}
	cache.RUnlock()

	for _, host := range Hosts.getHost() {
		if host.Domain == nil {
			continue
		}
		if host.Domain.MatchString(hostname) {
			cache.Lock()
			cache.m[hostname] = &host
			cache.Unlock()
			hostDealer(&host, res, req)
			return
		}
	}

	cache.Lock()
	cache.m[hostname] = false
	cache.Unlock()
}

func nonsecure(res http.ResponseWriter, req *http.Request) {
	req.Header.Del("X-Secure")
	proxy(res, req)
}

func secure(res http.ResponseWriter, req *http.Request) {
	req.Header.Set("X-Secure", "Monkey")
	proxy(res, req)
}

func Serve() {
	http.ListenAndServe(DefaultAddr, http.HandlerFunc(nonsecure))
}

func ServeTLS(certFile, keyFile string) {
	http.ListenAndServeTLS(DefaultAddrTLS, certFile, keyFile, http.HandlerFunc(secure))
}
