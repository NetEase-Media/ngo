package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/db"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "db"

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
	opt := db.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	client, err := db.New(opt)
	if err != nil {
		return err
	}
	db.SetClient(opt.Name, client)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		client.Close()
		return nil
	})
	return err
}
