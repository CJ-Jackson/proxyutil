package proxyutil

import (
	"net/http"
	"sort"
	"sync"
)

type balencerItem struct {
	handler http.Handler
	tally   uint8
}

type balencerItems []*balencerItem

func (bal balencerItems) Len() int {
	return len(bal)
}

func (bal balencerItems) Less(i, j int) bool {
	return bal[i].tally < bal[j].tally
}

func (bal balencerItems) Swap(i, j int) {
	bal[i], bal[j] = bal[j], bal[i]
}

func (bal balencerItems) Reset() {
	for _, ba := range bal {
		ba.tally = uint8(0)
	}
}

type balencer struct {
	sync.Mutex
	items balencerItems
}

func Balencer(handlers ...http.Handler) http.Handler {
	bal := balencer{}
	bal.Lock()
	defer bal.Unlock()
	bal.items = balencerItems{}

	for _, handler := range handlers {
		bal.items = append(bal.items, &balencerItem{handler, uint8(0)})
	}

	return bal
}

func (bal balencer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	bal.Lock()
	item := bal.items[0]
	item.tally++
	sort.Sort(bal.items)
	if item.tally >= 200 {
		bal.items.Reset()
	}
	bal.Unlock()
	item.handler.ServeHTTP(res, req)
}
