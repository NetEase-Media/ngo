package db

import (
	"sync"
)

var (
	mu        sync.RWMutex
	dbClients = make(map[string]*Client)
)

func SetClient(name string, client *Client) {
	mu.Lock()
	defer mu.Unlock()
	dbClients[name] = client
}

func GetClient(name string) *Client {
	mu.RLock()
	defer mu.RUnlock()
	return dbClients[name]
}
