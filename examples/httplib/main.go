package main

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/httplib"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		var rs string
		code, err := httplib.Get("https://www.163.com").BindString(&rs).Do(context.Background())
		if err != nil {
			log.Errorf("err: %s", err)
		}
		log.Infof("code: %d", code)
		return nil
	}
	app.Start()
}
