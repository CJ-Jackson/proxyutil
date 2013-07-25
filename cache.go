package proxyutil

import (
	"sync"
)

var cache = struct {
	sync.RWMutex
	m map[string]interface{}
}{
	m: map[string]interface{}{},
}
