package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/zookeeper"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "zookeeper"

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	cs := config.SubSlice(Key)
	if len(cs) == 0 {
		return nil
	}
	for i := range cs {
		if err := newFromConfig(cs[i]); err != nil {
			return err
		}
	}
	return nil
}

func newFromConfig(c *config.Configuration) error {
	opt := zookeeper.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	proxy, err := zookeeper.New(opt)
	if err != nil {
		return err
	}
	zookeeper.SetClient(opt.Name, proxy)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		proxy.Close()
		return nil
	})
	return nil
}
