package redis

import "github.com/go-redis/redis/v8"

func newClientOptions(opt *Options) *redis.Options {
	return &redis.Options{
		Addr:               opt.Addr[0],
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

func NewClient(opt *Options) *RedisContainer {
	baseClient := redis.NewClient(newClientOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeClient,
	}
	baseClient.AddHook(newMetricHook(c))
	baseClient.AddHook(newTracingHook(c))
	return c
}
