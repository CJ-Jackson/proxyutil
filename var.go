package proxyutil

import (
	"net/http"
)

var (
	Hosts          = &ListOfHosts{}
	DefaultAddr    = ":80"
	DefaultAddrTLS = ":443"

	Prepare http.Handler
	Finish  http.Handler
)
