// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memcache

import (
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	memcacheClients map[string]*MemcacheProxy // 用于保存所有配置的客户端
)

// Options 可配置的数据
type Options struct {
	Name         string        // 客户端名称，需要唯一
	Timeout      time.Duration // 客户端连接超时时间
	MaxIdleConns int           // 最大空闲连接
	Addr         []string      // 集群地址
}

// MemcacheProxy memcache 三方包的包装器类
type MemcacheProxy struct {
	base *memcache.Client
}

// NewMemcacheProxy 根据配置得到客户端
func NewMemcacheProxy(opt *Options) *MemcacheProxy {
	c := memcache.New(opt.Addr...)
	c.Timeout = opt.Timeout
	c.MaxIdleConns = opt.MaxIdleConns
	p := &MemcacheProxy{
		base: c,
	}
	return p
}

func Init(opts []Options) (err error) {
	if memcacheClients != nil {
		panic("memcache client has initialed, you should not initial it again!")
	}

	// 没有配置，则跳过初始化
	if len(opts) == 0 {
		return
	}
	memcacheClients, err = newClientFromOption(opts)
	return
}

func GetAllClients() []*MemcacheProxy {
	clients := make([]*MemcacheProxy, 0, len(memcacheClients))
	for _, v := range memcacheClients {
		clients = append(clients, v)
	}
	return clients
}

// GetClient
func GetClient(name string) (ret *MemcacheProxy) {
	return memcacheClients[name]
}

// 解析，构建客户端
func newClientFromOption(opts []Options) (map[string]*MemcacheProxy, error) {
	cs := make(map[string]*MemcacheProxy)

	for i := range opts {
		opt := &opts[i]

		if _, ok := cs[opt.Name]; ok {
			return nil, fmt.Errorf("duplicated memcache config %s", opt.Name)
		}

		if len(opt.Addr) == 0 {
			return nil, errors.New("memcache addr can not be empty!")
		}
		cs[opt.Name] = NewMemcacheProxy(opt)
	}
	return cs, nil
}
