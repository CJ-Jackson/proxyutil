Proxy Util
==========

Just a small collection of proxy utilites!

	package main

	import (
	  "github.com/CJ-Jackson/proxyutil"
	  "net/http"
	  "runtime"
	)

	func init() {
	  proxyutil.Prepare = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	    res.Header().Set("Server", "Bespoken")
	  })

	  proxyutil.Hosts.Register(
	    // cj-jackson.com, port 39001
	    proxyutil.Host{
	      Domain: proxyutil.Names("cj-jackson.com", "www.cj-jackson.com"),
	      Proxy:  proxyutil.ReverseProxy("http://127.0.0.1:39001"),
	      WS:     proxyutil.WebsocketProxy("127.0.0.1:39001"),
	    },
	  )
	}

	func main() {
	  runtime.GOMAXPROCS(runtime.NumCPU())
	  go proxyutil.Serve()
	  proxyutil.ServeTLS("./cert.pem", "./key.pem")
	}