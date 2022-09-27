package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/log"
)

const (
	Key = "log"
)

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	cs := config.SubSlice(Key)
	if len(cs) == 0 {
		one := config.Sub(Key)
		if one != nil {
			if err := newFromConfig(one); err != nil {
				return err
			}
		}
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
	opt := log.NewDefaultOptions()
	if err := c.Unmarshal(&opt); err != nil {
		return err
	}
	logger, err := log.New(opt)
	if err != nil {
		return err
	}
	log.SetLogger(opt.Name, logger)
	if opt.Name == log.DefaultLoggerName {
		log.SetDefaultLogger(logger)
	}
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		logger.Sync()
		return nil
	})
	return nil
}
