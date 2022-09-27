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
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/dlock/redis"
	"github.com/hashicorp/go-multierror"
)

// A DelayFunc is used to decide the amount of time to wait between retries.
type DelayFunc func(tries int) time.Duration

type Action func()

// A Mutex is a distributed mutual exclusion lock.
type Mutex struct {
	name   string
	expiry time.Duration

	tries     int
	delayFunc DelayFunc

	factor float64

	quorum int

	genValueFunc func() (string, error)
	value        string
	until        time.Time

	pools []redis.Pool

	action Action
}

// Name returns mutex name (i.e. the Redis key).
func (m *Mutex) Name() string {
	return m.name
}

// Value returns the current random value. The value will be empty until a lock is acquired (or WithValue option is used).
func (m *Mutex) Value() string {
	return m.value
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() (bool, error) {
	return m.LockContext(context.TODO())
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) LockContext(ctx context.Context) (bool, error) {
	value, err := m.genValueFunc()
	if err != nil {
		return false, err
	}

	for i := 0; i < m.tries; i++ {
		if i != 0 {
			time.Sleep(m.delayFunc(i))
		}

		start := time.Now()

		n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
			return m.acquire(ctx, pool, value)
		})
		if n == 0 && err != nil {
			return false, err
		}

		now := time.Now()
		until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.factor)))
		if n >= m.quorum && now.Before(until) {
			m.value = value
			m.until = until
			return true, nil
		}
		_, _ = m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
			return m.release(ctx, pool, value)
		})
		if i == m.tries-1 && err != nil {
			return false, err
		}
	}

	return false, nil
}

// Unlock unlocks m and returns the status of unlock.
func (m *Mutex) Unlock() (bool, error) {
	return m.UnlockContext(context.TODO())
}

// Unlock unlocks m and returns the status of unlock.
func (m *Mutex) UnlockContext(ctx context.Context) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.release(ctx, pool, m.value)
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *Mutex) Extend() (bool, error) {
	return m.ExtendContext(context.TODO())
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *Mutex) ExtendContext(ctx context.Context) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.touch(ctx, pool, m.value, int(m.expiry/time.Millisecond))
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

func (m *Mutex) Valid() (bool, error) {
	return m.ValidContext(context.TODO())
}

func (m *Mutex) ValidContext(ctx context.Context) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.valid(ctx, pool)
	})
	return n >= m.quorum, err
}

func (m *Mutex) DoContext(ctx context.Context) (bool, bool, error) {
	var executed bool
	succ, err := m.LockContext(ctx)
	if !succ {
		return succ, executed, err
	}

	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
		succ, err := m.UnlockContext(ctx)
		if !succ {
			log.Errorf("failed to unlock. name: %s, err: %+v", m.name, err)
		}
	}()
	childCtx, cancel := context.WithCancel(ctx)
	m.renew(childCtx)
	m.action()
	executed = true

	return succ, executed, err
}

func (m *Mutex) Do() (bool, bool, error) {
	return m.DoContext(context.TODO())
}

// WithExpiry can be used to set the expiry of a mutex to the given value.
func (m *Mutex) WithExpiry(expiry time.Duration) *Mutex {
	m.expiry = expiry
	return m
}

// WithTries can be used to set the number of times lock acquire is attempted.
func (m *Mutex) WithTries(tries int) *Mutex {
	m.tries = tries
	return m
}

// WithRetryDelay can be used to set the amount of time to wait between retries.
func (m *Mutex) WithRetryDelay(delay time.Duration) *Mutex {
	m.delayFunc = func(tries int) time.Duration {
		return delay
	}
	return m
}

// WithRetryDelayFunc can be used to override default delay behavior.
func (m *Mutex) WithRetryDelayFunc(delayFunc DelayFunc) *Mutex {
	m.delayFunc = delayFunc
	return m
}

// WithDriftFactor can be used to set the clock drift factor.
func (m *Mutex) WithDriftFactor(factor float64) *Mutex {
	m.factor = factor
	return m
}

// WithGenValueFunc can be used to set the custom value generator.
func (m *Mutex) WithGenValueFunc(genValueFunc func() (string, error)) *Mutex {
	m.genValueFunc = genValueFunc
	return m
}

// WithValue can be used to assign the random value without having to call lock. This allows the ownership of a lock to be "transfered" and allows the lock to be unlocked from elsewhere.
func (m *Mutex) WithValue(v string) *Mutex {
	m.value = v
	return m
}

func (m *Mutex) renew(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.expiry / 3):
				succ, err := m.ExtendContext(ctx)
				if !succ {
					log.Errorf("failed to renew. name: %s, err: %+v", m.name, err)
				} else {
					log.Infof("renew. name: %s success", m.name)
				}
			}
		}
	}()

}

func (m *Mutex) valid(ctx context.Context, pool redis.Pool) (bool, error) {
	if m.value == "" {
		return false, nil
	}
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	reply, err := conn.Get(m.name)
	if err != nil {
		return false, err
	}
	return m.value == reply, nil
}

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (m *Mutex) acquire(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	reply, err := conn.SetNX(m.name, value, m.expiry)
	if err != nil {
		return false, err
	}
	return reply, nil
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (m *Mutex) release(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	status, err := conn.Eval(deleteScript, m.name, value)
	if err != nil {
		return false, err
	}
	return status != int64(0), nil
}

var touchScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

func (m *Mutex) touch(ctx context.Context, pool redis.Pool, value string, expiry int) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	status, err := conn.Eval(touchScript, m.name, value, expiry)
	if err != nil {
		return false, err
	}
	return status != int64(0), nil
}

func (m *Mutex) actOnPoolsAsync(actFn func(redis.Pool) (bool, error)) (int, error) {
	type result struct {
		Status bool
		Err    error
	}

	ch := make(chan result)
	for _, pool := range m.pools {
		go func(pool redis.Pool) {
			r := result{}
			r.Status, r.Err = actFn(pool)
			ch <- r
		}(pool)
	}
	n := 0
	var err error
	for range m.pools {
		r := <-ch
		if r.Status {
			n++
		} else if r.Err != nil {
			err = multierror.Append(err, r.Err)
		}
	}
	return n, err
}
