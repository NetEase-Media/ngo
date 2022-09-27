package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/NetEase-Media/ngo/pkg/client/db"
	"github.com/NetEase-Media/ngo/pkg/client/httplib"
	"github.com/NetEase-Media/ngo/pkg/client/redis"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/NetEase-Media/ngo/pkg/tracing"
	"github.com/NetEase-Media/ngo/pkg/tracing/pinpoint"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	// 自定义 agent id
	pinpoint.SetAgentId("1.1.1.1")
	app := ngo.Init()
	s := http.Get()
	s.AddRoute(http.GET, "/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var t test
		db.WithContext(ctx, db.GetClient("db01")).
			Raw("select blackword from blacklist t where bid = ?", 1).Find(&t)
		log.Info(t.Blackword)
		client := redis.GetClient("redis01")
		client.Set(ctx, "key", "value", time.Second*5).Result()

		httplib.Get("http://localhost:8080/tracing").Do(ctx)

		c.String(http.StatusOK, "hello world! traceId: %s", tracing.GetTraceId(ctx))
	})

	s.AddRoute(http.GET, "/tracing", func(c *gin.Context) {
		c.String(http.StatusOK, "tracing")
	})

	go func() {
		time.Sleep(time.Second * 3)
		for i := 0; i < 10; i++ {
			go func() {
				for {
					intn := rand.Intn(3000) + 1
					time.Sleep(time.Duration(intn) * time.Millisecond)
					httplib.Get("http://localhost:8080/").Do(context.Background())
				}
			}()
		}
	}()

	app.Start()
}

type test struct {
	Blackword string `json:"blackword"`
}
