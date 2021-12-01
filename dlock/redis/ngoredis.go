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
	"context"
	"strings"
	"time"

	"github.com/NetEase-Media/ngo/client/redis"
	redisv8 "github.com/go-redis/redis/v8"
)

type pool struct {
	delegate redis.Redis
}

func (p *pool) Get(ctx context.Context) (Conn, error) {
	return &conn{p.delegate, ctx}, nil
}

func NewPool(delegate redis.Redis) Pool {
	return &pool{delegate}
}

type conn struct {
	delegate redis.Redis
	ctx      context.Context
}

func (c *conn) Get(name string) (string, error) {
	value, err := c.delegate.Get(c.ctx, name).Result()
	return value, noErrNil(err)
}

func (c *conn) Set(name string, value string) (bool, error) {
	reply, err := c.delegate.Set(c.ctx, name, value, 0).Result()
	return reply == "OK", err
}

func (c *conn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	return c.delegate.SetNX(c.ctx, name, value, expiry).Result()
}

func (c *conn) PTTL(name string) (time.Duration, error) {
	return c.delegate.PTTL(c.ctx, name).Result()
}

func (c *conn) Eval(script *Script, keysAndArgs ...interface{}) (interface{}, error) {
	keys := make([]string, script.KeyCount)
	args := keysAndArgs

	if script.KeyCount > 0 {
		for i := 0; i < script.KeyCount; i++ {
			keys[i] = keysAndArgs[i].(string)
		}
		args = keysAndArgs[script.KeyCount:]
	}

	v, err := c.delegate.EvalSha(c.ctx, script.Hash, keys, args...).Result()
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		v, err = c.delegate.Eval(c.ctx, script.Src, keys, args...).Result()
	}
	return v, noErrNil(err)
}

func (c *conn) Close() error {
	return nil
}

func (c *conn) _context(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return nil
}

func noErrNil(err error) error {
	if err == redisv8.Nil {
		return nil
	}
	return err
}
