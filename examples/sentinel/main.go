package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/sentinel"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		for i := 0; i < 100; i++ {
			e, b := sentinel.Entry("abc")
			if b != nil {
				log.Info("too many requests")
			} else {
				log.Info("pass")
				e.Exit()
			}
		}
		return nil
	}
	app.Start()
}
