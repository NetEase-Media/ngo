package kafka

import (
	"sync"
)

var (
	mu          sync.RWMutex
	consumerMap = make(map[string]*Consumer)
	producerMap = make(map[string]*Producer)
)

func GetConsumer(name string) *Consumer {
	mu.RLock()
	defer mu.RUnlock()
	return consumerMap[name]
}

func GetProducer(name string) *Producer {
	mu.RLock()
	defer mu.RUnlock()
	return producerMap[name]
}

func SetConsumer(name string, consumer *Consumer) {
	mu.Lock()
	defer mu.Unlock()
	consumerMap[name] = consumer
}

func SetProducer(name string, producer *Producer) {
	mu.Lock()
	defer mu.Unlock()
	producerMap[name] = producer
}
