package memcache

import (
	"errors"
	"time"
)

// Options 可配置的数据
type Options struct {
	Name         string        // 客户端名称，需要唯一
	Timeout      time.Duration // 客户端连接超时时间
	MaxIdleConns int           // 最大空闲连接
	Addr         []string      // 集群地址
}

func NewDefaultOptions() *Options {
	return &Options{}
}

func checkOptions(opt *Options) error {
	if opt.Name == "" {
		return errors.New("client name can not be nil")
	}
	return nil
}
