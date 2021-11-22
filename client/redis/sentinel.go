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

import "github.com/go-redis/redis/v8"

func newSentinelOptions(opt *Options) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:         opt.MasterNames[0],
		SentinelAddrs:      opt.Addr,
		Username:           opt.Username,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
	}
}

func NewSentinelClient(opt *Options) *redisContainer {
	baseClient := redis.NewFailoverClient(newSentinelOptions(opt))
	c := &redisContainer{
		Redis:     baseClient,
		opt:       *opt,
		redisType: RedisTypeSentinel,
	}
	baseClient.AddHook(newMetricHook(c))
	baseClient.AddHook(newTracingHook(c))
	return c
}
