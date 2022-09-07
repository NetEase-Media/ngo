package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/memcache"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "memcache"

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
	opt := memcache.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	proxy, err := memcache.New(opt)
	if err != nil {
		return err
	}
	memcache.SetClient(opt.Name, proxy)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		return nil
	})
	return nil
}
