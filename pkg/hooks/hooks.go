package hooks

import (
	"context"
	"sync"
)

type Stage int

const (
	Init Stage = iota
	ServerStart
	ComponentStop
	ServerStop
	ServerGracefulStop
)

var (
	globalHooks = make(map[Stage][]func(ctx context.Context) error)
	mu          = sync.RWMutex{}
)

func Register(stage Stage, fns ...func(ctx context.Context) error) {
	mu.Lock()
	defer mu.Unlock()
	globalHooks[stage] = append(globalHooks[stage], fns...)
}

func GetFns(stage Stage) []func(ctx context.Context) error {
	return globalHooks[stage]
}
