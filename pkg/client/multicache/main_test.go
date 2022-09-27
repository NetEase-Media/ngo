package multicache

import (
	"fmt"
	"os"
	"testing"

	"github.com/NetEase-Media/ngo/pkg/client/redis"
)

func TestMain(m *testing.M) {
	fmt.Printf("multicache start.")
	// cache := gcache.New(100000).Simple().Build()
	// InitLocal(cache)

	redisOpts := redis.Options{
		Name:     "client1",
		Addr:     []string{"127.0.0.1:6379"},
		Password: "rntestncr",
		ConnType: "client",
	}

	r, _ := redis.New(&redisOpts)
	redis.SetClient("client1", r)
	// redisClient := redis.GetClient("client1")
	// InitRedis(redisClient)

	opts := []Options{
		{
			Type:     "local",
			Priority: 0,
			Capacity: 100000,
		},
		{
			Type:         "redis",
			Priority:     1,
			DefaultRedis: "client1",
		},
	}
	Init(opts)
	m.Run()
	fmt.Printf("multicache end.")
	os.Exit(0)
}
