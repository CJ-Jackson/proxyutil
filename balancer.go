package proxyutil

import (
	"net/http"
	"sync"
)

type balencer struct {
	sync.Mutex
	handlers []http.Handler
}

func Balencer(handlers ...http.Handler) http.Handler {
	bal := &balencer{}
	bal.Lock()
	defer bal.Unlock()
	bal.handlers = []http.Handler{}

	for _, handler := range handlers {
		bal.handlers = append(bal.handlers, handler)
	}

	return bal
}

func (bal *balencer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	bal.Lock()
	if len(bal.handlers) < 0 {
		bal.Unlock()
		return
	}
	item := bal.handlers[0]
	if len(bal.handlers) <= 1 {
		goto unlock
	}
	bal.handlers = append(bal.handlers[1:], item)
unlock:
	bal.Unlock()
	item.ServeHTTP(res, req)
}
