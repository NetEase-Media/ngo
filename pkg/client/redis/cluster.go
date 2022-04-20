package redis

import "github.com/go-redis/redis/v8"

func newClusterOptions(opt *Options) *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:              opt.Addr,
		Username:           opt.Username,
		Password:           opt.Password,
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

func NewClusterClient(opt *Options) *RedisContainer {
	baseClient := redis.NewClusterClient(newClusterOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeCluster,
	}
	return c
}
