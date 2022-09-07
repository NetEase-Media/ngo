package zookeeper

import (
	"sync"

	"github.com/go-zookeeper/zk"
)

func New(opt *Options) (*ZookeeperProxy, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	proxy, err := NewZookeeperProxy(opt)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}

type ZookeeperProxy struct {
	Opt *Options

	Conn       *zk.Conn
	tmpNode    sync.Map
	listenCh   chan Event
	listenSign uint32
	stop       chan struct{}
	done       *sync.WaitGroup
}
