package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/dlock"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
)

const Key = "dlock"

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}

	opt := dlock.NewDefaultOptions()
	if err := config.UnmarshalKey(Key, &opt); err != nil {
		return err
	}
	d, err := dlock.New(opt)
	if err != nil {
		return err
	}
	dlock.SetDefaultDlock(d)
	return err
}
