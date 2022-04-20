package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/job/xxljob"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/service"
)

const (
	Key = "xxljob"
)

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}

	var xxljobOptions xxljob.Options
	if err := config.UnmarshalKey(Key, &xxljobOptions); err != nil {
		return err
	}
	if err := xxljob.Init(&xxljobOptions, service.GetClusterName()); err != nil {
		return err
	}
	hooks.Register(hooks.ComponentStop, func(ctx context.Context) error {
		xxljob.Stop()
		return nil
	})
	return nil
}
