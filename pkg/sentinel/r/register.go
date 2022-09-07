package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/sentinel"
)

const (
	Key = "sentinel"
)

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}

	var sentinelOptions sentinel.Options
	if err := config.UnmarshalKey(Key, &sentinelOptions); err != nil {
		return err
	}
	return sentinel.Init(&sentinelOptions)
}
