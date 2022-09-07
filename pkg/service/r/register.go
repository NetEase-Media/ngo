package r

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/service"
)

const Key = "service"

func init() {
	hooks.Register(hooks.Init, NewFromConfig)
}

func NewFromConfig(ctx context.Context) error {
	if !config.Exists(Key) {
		return nil
	}

	opt := service.NewDefaultOptions()
	if err := config.UnmarshalKey(Key, &opt); err != nil {
		return err
	}
	s, err := service.New(opt)
	if err != nil {
		return err
	}
	service.SetDefaultService(s)
	return err
}
