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

package redis

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/NetEase-Media/ngo/adapter/log"
)

const (
	RedisTypeClient          = "client"
	RedisTypeCluster         = "cluster"
	RedisTypeSentinel        = "sentinel"
	RedisTypeShardedSentinel = "sharded_sentinel"
)

type Redis interface {
	redis.Cmdable
	io.Closer
}

var (
	// redisClients 将redis client和cluster都集中存储起来
	redisClients map[string]*redisContainer
)

// Options 是redis客户端的通用配置选项，兼容单实例和cluster类型。
type Options struct {
	// 用户需要保证名字唯一
	Name string

	// 连接类型，必须指定。包含client、cluster、sentinel、sharded_sentinel四种类型。
	ConnType string

	// 地址列表，格式为host:port。如果是单实例只会取第一个。
	Addr []string

	// master 名称，只当sentinel、sharded_sentinel 类型必填。如果是sentinel只会取第一个。
	MasterNames []string

	// 自动生成分片名称，如果为false，默认使用MasterName， 只当sharded_sentinel 类型使用。
	// 该字段用来兼容旧项目，非特殊情况请勿设置成true，否则在MasterNames顺序变化时会造成分配rehash
	AutoGenShardName bool

	// 用于认证的用户名
	Username string

	// 用于认证的密码
	Password string

	// 所使用的数据库
	DB int

	// 最大重试次数
	MaxRetries int

	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	// 超时时间
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// 最大连接数
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	// TODO: 未来增加
	TLSConfig *tls.Config
}

func NewDefaultOptions() *Options {
	return &Options{}
}

func Init(opts []Options) (err error) {
	if redisClients != nil {
		panic("duplicated redis init")
	}

	if len(opts) == 0 {
		// 如果没有redis配置则跳过
		log.Info("empty redis config, so skip init")
		return
	}

	redisClients, err = newFromConfig(opts)
	return
}

// newFromConfig 解析yaml配置，初始化客户端
func newFromConfig(opts []Options) (m map[string]*redisContainer, err error) {
	m = make(map[string]*redisContainer)
	for i := range opts {
		opt := &opts[i]

		// 重复配置判断
		if _, ok := m[opt.Name]; ok {
			return nil, fmt.Errorf("duplicated redis config %s", opt.Name)
		}

		// 检查地址合法性
		if len(opt.Addr) == 0 {
			return nil, errors.New("empty address")
		}

		var c *redisContainer
		// 判断连接类型
		switch opt.ConnType {
		case RedisTypeClient:
			c = NewClient(opt)
		case RedisTypeCluster:
			c = NewClusterClient(opt)
		case RedisTypeSentinel:
			if len(opt.MasterNames) == 0 {
				err = errors.New("empty master name")
				return
			}
			c = NewSentinelClient(opt)
		case RedisTypeShardedSentinel:
			if len(opt.MasterNames) == 0 {
				err = errors.New("empty master name")
				return
			}
			c = NewShardedSentinelClient(opt)
		default:
			err = errors.New("redis connection type need ")
			return
		}

		m[c.opt.Name] = c
	}
	return
}

func GetClient(name string) Redis {
	redis, ok := redisClients[name]
	if ok {
		return redis
	}
	return nil
}

func StopAll() {
	for name, client := range redisClients {
		client.Close()
		log.Infof("Stop redis client %s", name)
	}
}

// redisContainer 用来存储redis客户端及其额外信息
type redisContainer struct {
	Redis
	opt       Options
	redisType string
}
