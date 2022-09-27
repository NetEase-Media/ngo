package multicache

import (
	"context"
	"time"

	"github.com/NetEase-Media/ngo/pkg/client/redis"
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
