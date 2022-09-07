package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/redis"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "redis"

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
	opt := redis.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	r, err := redis.New(opt)
	if err != nil {
		return err
	}
	redis.SetClient(opt.Name, r)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		r.Close()
		return nil
	})
	return nil
}
