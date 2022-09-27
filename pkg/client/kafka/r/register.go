package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/kafka"
	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "kafka"

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
	opt := kafka.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	k, err := kafka.New(opt)
	if err != nil {
		return err
	}
	kafka.SetProducer(opt.Name, k.Producer)
	kafka.SetConsumer(opt.Name, k.Consumer)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		k.Consumer.Stop()
		k.Producer.Close()
		return nil
	})
	return nil
}
