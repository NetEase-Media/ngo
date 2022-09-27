package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/server/http"
)

const Key = "httpServer"

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}

	sub := config.Sub(Key)
	if sub == nil {
		sub = config.Empty()
	}
	return newFromConfig(sub)
}

func newFromConfig(c *config.Configuration) error {
	opt := http.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	server, err := http.New(opt)
	if err != nil {
		return err
	}
	http.Set(server)
	hooks.Register(hooks.ServerStart, func(ctx context.Context) error {
		return server.Start()
	})
	hooks.Register(hooks.ServerStop, func(ctx context.Context) error {
		return server.Stop()
	})
	hooks.Register(hooks.ServerGracefulStop, func(ctx context.Context) error {
		return server.GracefulStop(ctx)
	})
	return nil
}
