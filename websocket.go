package proxyutil

import (
	"io"
	"log"
	"net"
	"net/http"
)

func WebsocketProxy(target string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)
		r.Header.Set("X-Forwarded-For", clientIp)

		d, err := net.Dial("tcp", target)
		if err != nil {
			http.Error(w, "Error contacting backend server.", 500)
			log.Printf("Error dialing websocket backend %s: %v", target, err)
			return
		}
		defer d.Close()
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Not a hijacker?", 500)
			return
		}
		nc, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Hijack error: %v", err)
			return
		}
		defer nc.Close()

		err = r.Write(d)
		if err != nil {
			log.Printf("Error copying request to target: %v", err)
			return
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(d, nc)
		go cp(nc, d)
		<-errc
	})
}

func IsWebSocket(req *http.Request) bool {
	return req.Header.Get("Connection") == "Upgrade" && req.Header.Get("Upgrade") == "websocket"
}
