package etcd

import (
	"context"
	"time"

	"github.com/NetEase-Media/ngo/pkg/util"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Option struct {
	Endpoints      []string `form:"endpoints"`
	UserName       string   `form:"user_name"`
	Password       string   `form:"password"`
	ReadTimeoutSec int      `form:"read_timeout_sec"`
	Key            string   `form:"key"`
	WithPrefix     bool     `form:"with_prefix"`
}

type Client struct {
	Client *clientv3.Client
	Opt    Option
	stop   chan struct{}
}

// NewDataSource 新建 etcd client
func NewDataSource(opt Option) (*Client, error) {
	if opt.ReadTimeoutSec == 0 {
		opt.ReadTimeoutSec = 3
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints: opt.Endpoints,
		Username:  opt.UserName,
		Password:  opt.Password,
	})
	if err != nil {
		return nil, err
	}
	return &Client{Client: client, Opt: opt, stop: make(chan struct{})}, nil
}

func (c *Client) ReadConfig() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Opt.ReadTimeoutSec))
	defer cancel()
	resp, err := c.Client.Get(ctx, c.Opt.Key)
	if err != nil {
		return nil, err
	}

	return resp.Kvs[0].Value, nil
}

func (c *Client) IsConfigChanged() <-chan struct{} {
	notify := make(chan struct{})

	opts := make([]clientv3.OpOption, 0)
	if c.Opt.WithPrefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	ch := c.Client.Watch(context.Background(), c.Opt.Key, opts...)

	util.GoWithRecover(func() {
		for {
			select {
			case <-c.stop:
				return
			case <-ch:
				notify <- struct{}{}
			}
		}
	}, nil)
	return notify
}

func (c *Client) Close() error {
	c.stop <- struct{}{}
	return nil
}
