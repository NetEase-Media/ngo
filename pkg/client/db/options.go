package db

import (
	"errors"
	"time"
)

// Options 是MysqlClient的配置数据
type Options struct {
	Name            string
	Type            string
	Url             string
	MaxIdleCons     int
	MaxOpenCons     int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDefaultOptions() *Options {
	return &Options{
		Type:            "mysql",
		MaxIdleCons:     10,
		MaxOpenCons:     10,
		ConnMaxLifetime: time.Second * 1000,
		ConnMaxIdleTime: time.Second * 60,
	}
}

func checkOptions(opt *Options) error {
	if opt.Name == "" {
		return errors.New("client name can not be nil")
	}
	return nil
}
