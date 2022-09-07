package memcache

import (
	"sync"
)

var (
	mu              sync.RWMutex
	memcacheClients = make(map[string]*MemcacheProxy)
)

func SetClient(name string, client *MemcacheProxy) {
	mu.Lock()
	defer mu.Unlock()
	memcacheClients[name] = client
}

func GetClient(name string) *MemcacheProxy {
	mu.RLock()
	defer mu.RUnlock()
	return memcacheClients[name]
}
