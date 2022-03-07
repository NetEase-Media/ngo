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

package multicache

import (
	"errors"
	"sort"

	"github.com/NetEase-Media/ngo/client/redis"
	"github.com/bluele/gcache"
)

const (
	defaultCapacity int = 1000000 // 默认最大100w
)

var MinerNotInitError = errors.New("miner has not init yet.")

// 这里主要是进行客户端的初始化
type Options struct {
	Type         string // 类型 local, redis 目前先支持这两个
	Priority     int    // 加载优先级, 0 是最高优先级数字越大优先级越低
	Capacity     int    // 最大的容积, 主要是给本地缓存使使用
	DefaultRedis string // 默认获取哪个 Redis 客户端
	Strategy     string // 淘汰策略: simple/lru/lfu/arc, 这个只有localcache支持
	Onload       bool   // 是否支持onload
}

var minerClient *Miner

func Init(opts []Options) { // 默认的客户端的名字, 默认取第一个配置项
	if len(opts) == 0 {
		return
	}

	if minerClient != nil {
		panic("local cache miner has multi initialed")
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Priority > opts[j].Priority
	})

	var handlers []Handler
	for _, cli := range opts {
		if cli.Capacity == 0 {
			cli.Capacity = defaultCapacity
		}
		switch cli.Type {
		case "local":
			if cli.Strategy == "" {
				cli.Strategy = "simple"
			}
			builder := gcache.New(cli.Capacity).EvictType(cli.Strategy)
			if cli.Onload {
				if OnloadFunc == nil {
					panic("OnloadFunc need initial first.")
				}
				builder.LoaderFunc(OnloadFunc)
			}
			cac := builder.Build()
			InitLocal(cac)
			localClient := GetLocalMiner()
			handlers = append(handlers, &localClient)
		case "redis":
			if cli.DefaultRedis == "" {
				// 没有拿到默认的redis
				panic("you need config default redis for multicache first.")
			}
			red := redis.GetClient(cli.DefaultRedis)
			if red == nil {
				panic("default redis do not exist.")
			}
			InitRedis(red)
			localRedis := GetRedisMiner()
			handlers = append(handlers, &localRedis)
		}
	}
	minerClient = &Miner{}
	minerClient.RegisterHander(handlers...)
}

// GetMiner 获取多级缓存处理客户端
func GetMiner() (*Miner, error) {
	if minerClient == nil {
		return nil, MinerNotInitError
	}
	return minerClient, nil
}

// ***** 主要的外部接口均在这里实现 *****

// Miner 设计思想是读的时候，先读本地缓存，获取不到才去redis读取数据
// 但是写数据的时候，则优先写入redis, 如果写入失败则不再写入，写入成功才会进行本地缓存的写入
type Miner struct {
	handersR []Handler // 读顺序
	handersW []Handler // 写顺序
}

var NoReaderHanlderError = errors.New("No Reader Handler Error.")
var NoWriterHanlderError = errors.New("No Writer Handler Error.")

// RegisterHander 注册Handler
func (m *Miner) RegisterHander(handlers ...Handler) {

	if len(m.handersR) > 0 || len(m.handersW) > 0 || len(m.handersR) != len(m.handersW) {
		panic("register handler failed!")
	}

	if len(handlers) < 1 {
		return
	}

	// 读的顺序
	m.handersR = append(m.handersR, handlers...)
	sort.Slice(m.handersR, func(i, j int) bool {
		return m.handersR[i].Priority() > m.handersR[j].Priority()
	})

	// 写的顺序
	m.handersW = append(m.handersW, handlers...)
	sort.Slice(m.handersR, func(i, j int) bool {
		return m.handersR[i].Priority() < m.handersR[j].Priority()
	})
}

func (m *Miner) Set(key, value string) (bool, error) {
	if len(m.handersW) == 0 {
		return false, NoWriterHanlderError
	}
	var ret bool
	var err error
	for _, h := range m.handersW {
		ret, err = h.Set(key, value)
		if err != nil {
			return false, err
		}
	}
	return ret, nil
}

func (m *Miner) SetWithTimeout(key, value string, ttl int) (bool, error) {
	if len(m.handersW) == 0 {
		return false, NoWriterHanlderError
	}
	var ret bool
	var err error
	for _, h := range m.handersW {
		ret, err = h.SetWithTimeout(key, value, ttl)
		if err != nil {
			return false, err
		}
	}
	return ret, nil
}

func (m *Miner) Get(key string) (string, error) {
	if len(m.handersR) == 0 {
		return "", NoReaderHanlderError
	}
	var ret string
	var err error
	var index int
	for _, h := range m.handersR {
		ret, err = h.Get(key)
		if err != nil {
			return "", err
		}
		if ret != "" { // 找到就提前返回
			if index != 0 {
				go m.sync(key, ret, index)
			}
			return ret, nil
		}
		index += 1 // 代表在非底层找到, 访问了一次
	}
	return ret, nil
}

// 底层缓存失效后，将会同步新的数据到底层缓存, 默认延迟60s
func (m *Miner) sync(key, value string, index int) {
	for i := 0; i < index; i++ {
		m.handersR[i].SetWithTimeout(key, value, 60)
	}
}

func (m *Miner) Evict(key string) (bool, error) {
	if len(m.handersW) == 0 {
		return false, NoWriterHanlderError
	}
	var ret bool
	var err error
	for _, h := range m.handersW {
		ret, err = h.Evict(key)
		if err != nil {
			return false, err
		}
	}
	return ret, nil
}

func (m *Miner) Clear() (bool, error) {
	if len(m.handersW) == 0 {
		return false, NoWriterHanlderError
	}
	var ret bool
	var err error
	for _, h := range m.handersW {
		ret, err = h.Clear()
		if err != nil {
			return false, err
		}
	}
	return ret, nil
}
