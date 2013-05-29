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

type StrDomain string

func (s StrDomain) MatchString(str string) bool {
	return string(s) == str
}

type StrDomains []string

func (s StrDomains) MatchString(str string) bool {
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
		return StrDomain(str[0])
	}
	return StrDomains(str)
}

func RegExp(pattern string) *regexp.Regexp {
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

var Hosts = &ListOfHosts{}

func hostnameWithoutPort(str string) string {
	hostname, _, err := net.SplitHostPort(str)
	if err != nil {
		hostname = strings.Split(str, ":")[0]
	}
	return hostname
}

func proxy(res http.ResponseWriter, req *http.Request) {
	hostname := hostnameWithoutPort(req.Host)
	for _, host := range Hosts.getHost() {
		if host.Domain.MatchString(hostname) {
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
			return
		}
	}
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
	http.ListenAndServe(":80", http.HandlerFunc(nonsecure))
}

func ServeTLS(certFile, keyFile string) {
	http.ListenAndServeTLS(":443", certFile, keyFile, http.HandlerFunc(secure))
}
