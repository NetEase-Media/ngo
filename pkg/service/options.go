package service

import (
	"errors"

	"github.com/NetEase-Media/ngo/pkg/env"
)

// Options
type Options struct {
	AppName     string
	ClusterName string
}

func NewDefaultOptions() *Options {
	return &Options{
		AppName:     env.GetAppName(),
		ClusterName: env.GetClusterName(),
	}
}

func checkOptions(o *Options) error {
	if o.AppName == "" || o.ClusterName == "" {
		return errors.New("lack of service config")
	}
	return nil
}
