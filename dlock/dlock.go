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

package dlock

import (
	"errors"
	"math/rand"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	ngoredis "github.com/NetEase-Media/ngo/client/redis"
	"github.com/NetEase-Media/ngo/dlock/redis"
)

const (
	minRetryDelayMilliSec = 50
	maxRetryDelayMilliSec = 250
)

var defaultDlock *Dlock

func Init(opt Options) error {
	if defaultDlock != nil {
		panic("duplicated init dlock")
	}
	if len(opt.Pools) == 0 {
		log.Info("empty dlock config, so skip init")
		return nil
	}

	pools := make([]redis.Pool, 0, 10)
	for _, v := range opt.Pools {
		rc := ngoredis.GetClient(v)
		if rc != nil {
			pool := redis.NewPool(rc)
			pools = append(pools, pool)
		}
	}
	if len(pools) < 1 {
		return errors.New("there must be at least one pool")
	}
	defaultDlock = New(pools...)
	return nil
}

func New(pools ...redis.Pool) *Dlock {
	return &Dlock{
		pools: pools,
	}
}

func NewMutex(key string, action Action) *Mutex {
	return defaultDlock.NewMutex(key, action)
}

type Options struct {
	Pools []string
}

type Dlock struct {
	pools []redis.Pool
}

func (d *Dlock) NewMutex(name string, action Action) *Mutex {
	return &Mutex{
		name:   name,
		expiry: 8 * time.Second,
		tries:  32,
		delayFunc: func(tries int) time.Duration {
			return time.Duration(rand.Intn(maxRetryDelayMilliSec-minRetryDelayMilliSec)+minRetryDelayMilliSec) * time.Millisecond
		},
		genValueFunc: genValue,
		factor:       0.01,
		quorum:       len(d.pools)/2 + 1,
		pools:        d.pools,
		action:       action,
	}
}
