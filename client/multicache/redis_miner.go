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
	"context"
	"time"

	"github.com/NetEase-Media/ngo/client/redis"
	redisV8 "github.com/go-redis/redis/v8"
)

// 本文件主要是处理reids相关的默认实现方式

// RedisMiner redis 相关具体实现
type RedisMiner struct {
	redis redis.Redis
}

var redisMiner RedisMiner

// InitRedisMiner
func InitRedis(redis redis.Redis) {
	redisMiner = RedisMiner{
		redis,
	}
}

func GetRedisMiner() RedisMiner {
	return redisMiner
}

func (r *RedisMiner) Priority() int {
	return 1
}

func (r *RedisMiner) Set(key, value string) (bool, error) {
	_, err := r.redis.Set(context.Background(), key, value, 0).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

// SetWithTimeout 时间单位为S
func (r *RedisMiner) SetWithTimeout(key, value string, ttl int) (bool, error) {
	_, err := r.redis.Set(context.Background(), key, value, time.Duration(int(time.Second))).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisMiner) Get(key string) (string, error) {
	ret, err := r.redis.Get(context.Background(), key).Result()
	if err == redisV8.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (r *RedisMiner) Evict(key string) (bool, error) {
	err := r.redis.Del(context.Background(), key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisMiner) Clear() (bool, error) {
	// 不支持清理所有数据，因此目前不做操作
	return true, nil
}
