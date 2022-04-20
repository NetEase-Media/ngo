package util

import (
	"context"
	"errors"
	"sync"
)

func NewWaitGroup() *WaitGroup {
	return new(WaitGroup)
}

func WithContext(ctx context.Context) (*WaitGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &WaitGroup{cancel: cancel}, ctx
}

type WaitGroup struct {
	cancel  func()
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func (g *WaitGroup) Run(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *WaitGroup) RunWithRecover(fn func() error) {
	g.wg.Add(1)

	GoWithRecover(func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}, func(r interface{}) {
		g.errOnce.Do(func() {
			g.err = errors.New("panic")
			if g.cancel != nil {
				g.cancel()
			}
		})
	})
}

func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
