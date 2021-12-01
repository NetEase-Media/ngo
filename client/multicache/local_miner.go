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
	"time"

	"github.com/bluele/gcache"
)

// 本文件主要是处理本地缓存相关的操作
var localMiner LocalMiner

// OnloadFunc 获取不到数据的时候支持加载数据
var OnloadFunc gcache.LoaderFunc

// InitOnloadFunc 需要应用自己初始化
func InitOnloadFunc(f gcache.LoaderFunc) {
	OnloadFunc = f
}

// LocalMiner 本地缓存的具体实现
type LocalMiner struct {
	cache gcache.Cache
}

// InitLocalMiner
func InitLocal(cache gcache.Cache) {
	localMiner = LocalMiner{
		cache,
	}
}

func GetLocalMiner() LocalMiner {
	return localMiner
}

func (l *LocalMiner) Priority() int {
	return 0
}

func (l *LocalMiner) Set(key, value string) (bool, error) {
	err := l.cache.Set(key, value)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SetWithTimeout 时间单位为S
func (l *LocalMiner) SetWithTimeout(key, value string, ttl int) (bool, error) {
	err := l.cache.SetWithExpire(key, value, time.Duration(ttl*int(time.Second)))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocalMiner) Get(key string) (string, error) {
	ret, err := l.cache.Get(key)
	if err == gcache.KeyNotFoundError { // key not found 不会当作异常处理
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ret.(string), nil
}

func (l *LocalMiner) Evict(key string) (bool, error) {
	l.cache.Remove(key)
	return true, nil
}

func (l *LocalMiner) Clear() (bool, error) {
	l.cache.Purge()
	return true, nil
}
