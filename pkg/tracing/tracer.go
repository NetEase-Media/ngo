package tracing

import (
	"fmt"
	"sync"
)

type Factory func(*Options) (Tracer, error)

var (
	mu            sync.RWMutex
	tracerFactory = make(map[string]Factory, 4)
)

func New(opt *Options) (Tracer, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	factory := GetFactory(opt.Type)
	if factory == nil {
		return nil, fmt.Errorf("unsupported type %s", opt.Type)
	}
	return factory(opt)
}

func Register(t string, f Factory) {
	mu.Lock()
	defer mu.Unlock()
	tracerFactory[t] = f
}

func GetFactory(t string) Factory {
	mu.RLock()
	defer mu.RUnlock()
	return tracerFactory[t]
}
