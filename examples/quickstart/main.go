package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	s := http.Get()
	s.AddRoute(http.GET, "/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world!")
	})
	app.Start()
}
