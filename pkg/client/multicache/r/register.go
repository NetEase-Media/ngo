package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/multicache"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "multicache"

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}
	//TODO: no default value
	var minerOptions []multicache.Options
	err := config.UnmarshalKey(Key, &minerOptions)
	if err != nil {
		return err
	}
	multicache.Init(minerOptions)
	return nil
}
