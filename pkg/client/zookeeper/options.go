package zookeeper

import (
	"errors"
	"time"
)

func NewDefaultOptions() *Options {
	return &Options{}
}

func checkOptions(opt *Options) error {
	if opt.Name == "" {
		return errors.New("client name can not be nil")
	}

	if len(opt.Addr) == 0 {
		return errors.New("zk: server list must not be empty")
	}

	if opt.SessionTimeout == 0 {
		opt.SessionTimeout = 5 * time.Second
	}
	return nil
}

type Options struct {
	Name           string        // 客户端名称， 需要唯一
	Addr           []string      // 节点地址
	SessionTimeout time.Duration // 连接创建超时时间
}
