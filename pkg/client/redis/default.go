package redis

import "sync"

var (
	mu           sync.RWMutex
	redisClients = make(map[string]Redis)
)

func SetClient(name string, client Redis) {
	mu.Lock()
	defer mu.Unlock()
	redisClients[name] = client
}

func GetClient(name string) Redis {
	mu.RLock()
	defer mu.RUnlock()
	return redisClients[name]
}
