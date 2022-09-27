package dlock

import (
	"errors"
	"math/rand"
	"time"

	ngoredis "github.com/NetEase-Media/ngo/pkg/client/redis"
	"github.com/NetEase-Media/ngo/pkg/dlock/redis"
)

const (
	minRetryDelayMilliSec = 50
	maxRetryDelayMilliSec = 250
)

var defaultDlock *Dlock

func SetDefaultDlock(dlock *Dlock) {
	defaultDlock = dlock
}

func New(opt *Options) (*Dlock, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	pools := make([]redis.Pool, 0, 10)
	for _, v := range opt.Pools {
		rc := ngoredis.GetClient(v)
		if rc != nil {
			pool := redis.NewPool(rc)
			pools = append(pools, pool)
		}
	}

	dlock, err := NewWithPools(pools...)
	if err != nil {
		return nil, err
	}
	return dlock, nil
}

func NewWithPools(pools ...redis.Pool) (*Dlock, error) {
	if len(pools) < 1 {
		return nil, errors.New("there must be at least one pool")
	}
	dlock := &Dlock{
		pools: pools,
	}
	return dlock, nil
}

func NewMutex(key string, action Action) *Mutex {
	return defaultDlock.NewMutex(key, action)
}

type Options struct {
	Pools []string
}

func NewDefaultOptions() *Options {
	return &Options{}
}

func checkOptions(opt *Options) error {
	if len(opt.Pools) == 0 {
		return errors.New("empty pools")
	}
	return nil
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
