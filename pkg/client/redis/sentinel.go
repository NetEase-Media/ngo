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

func NewSentinelClient(opt *Options) *RedisContainer {
	baseClient := redis.NewFailoverClient(newSentinelOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeSentinel,
	}
	baseClient.AddHook(newTracingHook(c))
	return c
}
