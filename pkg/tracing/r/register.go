package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/tracing"
)

const (
	Key = "tracing"
)

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}
	opt := tracing.NewDefaultOptions()
	if err := config.UnmarshalKey(Key, &opt); err != nil {
		return err
	}

	tracer, err := tracing.New(opt)
	if err != nil {
		return err
	}

	tracing.SetTracer(tracer)
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		tracer.Stop()
		return nil
	})

	config.OnChange(func(configuration *config.Configuration) {
		oldRate := opt.Pinpoint.Sampling.Rate
		newRate := configuration.GetInt("tracing.pinpoint.sampling.rate")
		if oldRate != newRate {
			log.Infof("tracing rate %v changed to %v", oldRate, newRate)
			tracer.SetSamplingRate(newRate)
			opt.Pinpoint.Sampling.Rate = newRate
		}
	})
	return nil
}
