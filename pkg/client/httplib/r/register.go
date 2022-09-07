package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/httplib"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const (
	Key = "httpClient"
)

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
	opt := httplib.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	client, err := httplib.New(opt)
	if err != nil {
		return err
	}
	httplib.SetDefaultHttpClient(client)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		client.Close()
		return nil
	})
	return nil
}
