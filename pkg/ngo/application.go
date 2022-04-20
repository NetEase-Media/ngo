package ngo

import (
	"context"
	"runtime"
	"sync"

	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/util"
)

func init() {
	// 设置procs的数量为容器cpu的2倍
	runtime.GOMAXPROCS(runtime.GOMAXPROCS(0) * 2)
}

func Init() *application {
	Flag()

	app := &application{}
	app.Init()
	return app
}

type application struct {
	Version   string
	PreStart  func() error
	AfterStop func() error
	cycle     *util.Cycle

	initOnce sync.Once
	stopOnce sync.Once
}

func (a *application) Init() {
	a.initOnce.Do(func() {
		ctx := context.Background()
		a.cycle = util.NewCycle()
		fns := hooks.GetFns(hooks.Init)
		for i := range fns {
			if err := fns[i](ctx); err != nil {
				util.CheckError(err)
			}
		}
	})
}

func (a *application) Start() {
	ctx := context.Background()
	if a.PreStart != nil {
		a.PreStart()
	}
	a.startServers(ctx)
	log.Info("shutdown")
}

func (a *application) GracefulStop() {
	a.stopOnce.Do(func() {
		ctx := context.Background()
		a.gracefulStopServers(ctx)
		a.stopComponents(ctx)
		if a.AfterStop != nil {
			a.AfterStop()
		}
		a.cycle.Close()
	})
}

func (a *application) Stop() {
	a.stopOnce.Do(func() {
		ctx := context.Background()
		a.stopServers(ctx)
		a.stopComponents(ctx)
		if a.AfterStop != nil {
			a.AfterStop()
		}
		a.cycle.Close()
	})
}

func (a *application) stopComponents(ctx context.Context) error {
	fns := hooks.GetFns(hooks.ComponentStop)
	for i := len(fns) - 1; i >= 0; i-- {
		if err := fns[i](ctx); err != nil {
			log.Errorl("stop component failed", log.String("error", err.Error()))
		}
	}
	return nil
}

func (a *application) startServers(ctx context.Context) {
	fns := hooks.GetFns(hooks.ServerStart)
	if len(fns) == 0 {
		return
	}
	for i := range fns {
		f := fns[i]
		a.cycle.Run(func() error {
			return f(ctx)
		})
	}
	if err := <-a.cycle.Wait(); err != nil {
		log.Errorl("start server failed", log.String("error", err.Error()))
	}
}

func (a *application) gracefulStopServers(ctx context.Context) {
	fns := hooks.GetFns(hooks.ServerGracefulStop)
	if len(fns) == 0 {
		return
	}
	for i := range fns {
		f := fns[i]
		a.cycle.Run(func() error {
			if err := f(ctx); err != nil {
				log.Errorl("stop server failed", log.String("error", err.Error()))
			}
			return nil
		})
	}
	<-a.cycle.Done()
}

func (a *application) stopServers(ctx context.Context) {
	fns := hooks.GetFns(hooks.ServerStop)
	if len(fns) == 0 {
		return
	}
	for i := range fns {
		f := fns[i]
		a.cycle.Run(func() error {
			if err := f(ctx); err != nil {
				log.Errorl("stop server failed", log.String("error", err.Error()))
			}
			return nil
		})
	}
	<-a.cycle.Done()
}
