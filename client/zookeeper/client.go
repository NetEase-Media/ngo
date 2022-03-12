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

package zookeeper

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	zookeeperClients map[string]*ZookeeperProxy // 用于保存所有配置的客户端
)

type Options struct {
	Name           string        // 客户端名称， 需要唯一
	Addr           []string      // 节点地址
	SessionTimeout time.Duration // 连接创建超时时间
}

type ZookeeperProxy struct {
	Conn *zk.Conn
}

func Init(opts []Options) (err error) {
	if zookeeperClients != nil {
		panic("zookeeper client has initialed, you should not initial it again!")
	}
	if len(opts) == 0 {
		return
	}

	zookeeperClients, err = NewClientFromOption(opts)
	return
}

func NewClientFromOption(opts []Options) (map[string]*ZookeeperProxy, error) {
	clients := make(map[string]*ZookeeperProxy)

	for i := range opts {
		opt := &opts[i]
		if opt.Name == "" {
			return nil, errors.New("client name can not be nil")
		}
		if _, ok := clients[opt.Name]; ok {
			return nil, fmt.Errorf("duplicated zookeeper config %s", opt.Name)
		}

		if len(opt.Addr) == 0 {
			return nil, errors.New("zk: server list must not be empty")
		}
		conn, _, err := zk.Connect(opt.Addr, opt.SessionTimeout)

		if err != nil {
			return nil, fmt.Errorf("connection failed")
		}
		clients[opt.Name] = &ZookeeperProxy{
			Conn: conn,
		}
	}
	return clients, nil
}
