package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/NetEase-Media/ngo/pkg/server/http/ratelimiter"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	s := http.Get()
	s.AddRoute(http.GET, "/", ratelimiter.RateLimiter(ratelimiter.WithResource("abc"),
		ratelimiter.WithErrorHandler(errorResponse)), doSomething)
	app.Start()
}

func errorResponse(c *gin.Context) {
	c.String(http.StatusTooManyRequests, "limit exceeded")
}

func doSomething(c *gin.Context) {
	c.String(http.StatusOK, "success")
}
