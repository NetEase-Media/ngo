package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/NetEase-Media/ngo/pkg/server/http/autoroute"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	s := http.Get()
	autoroute.AutoRoute(s, "/api/v1", NewHelloController())
	app.Start()
}
