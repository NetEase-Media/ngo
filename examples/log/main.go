package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		log.Info("info log")
		return nil
	}
	app.Start()
}
