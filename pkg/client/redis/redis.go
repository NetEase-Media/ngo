package redis

import (
	"errors"
	"io"

	"github.com/go-redis/redis/v8"
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

func New(opt *Options) (*RedisContainer, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}

	var c *RedisContainer
	// 判断连接类型
	switch opt.ConnType {
	case RedisTypeClient:
		c = NewClient(opt)
	case RedisTypeCluster:
		c = NewClusterClient(opt)
	case RedisTypeSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewSentinelClient(opt)
	case RedisTypeShardedSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewShardedSentinelClient(opt)
	default:
		err := errors.New("redis connection type need ")
		return nil, err
	}

	return c, nil
}

// RedisContainer 用来存储redis客户端及其额外信息
type RedisContainer struct {
	Redis
	Opt       Options
	redisType string
}
