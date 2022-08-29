package main

import (
	"context"
	"time"

	"github.com/NetEase-Media/ngo/pkg/client/redis"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

var (
	client redis.Redis
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()

	app.PreStart = func() error {
		key := "key"
		value := "value"

		client = redis.GetClient("redis01")
		_, err := client.Set(context.Background(), key, value, time.Second*5).Result()
		if err != nil {
			log.Error(err)
		}
		res, err := client.Get(context.Background(), key).Result()
		if err != nil {
			log.Error(err)
		}
		log.Info(res)
		return nil
	}
	app.Start()
}
