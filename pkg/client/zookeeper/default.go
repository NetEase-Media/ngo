package zookeeper

import (
	"sync"
)

var (
	mu               sync.RWMutex
	zookeeperClients = make(map[string]*ZookeeperProxy)
)

func SetClient(name string, client *ZookeeperProxy) {
	mu.Lock()
	defer mu.Unlock()
	zookeeperClients[name] = client
}

func GetClient(name string) *ZookeeperProxy {
	mu.RLock()
	defer mu.RUnlock()
	return zookeeperClients[name]
}
