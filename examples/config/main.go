package main

import (
	"github.com/NetEase-Media/ngo/pkg/config"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

//  go run . -c app.yaml
//  go run . -c app.properties
//  go run . -c apollo://106.54.227.205:8080?appId=ngo&cluster=ngo-demo&namespaceNames=application.properties,log.properties&configType=properties
//  go run . -c apollo://106.54.227.205:8080?appId=ngo&cluster=ngo-demo&namespaceNames=app.yaml,httpserver.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		test := config.GetString("service.appName")
		log.Info(test)
		return nil
	}
	app.Start()
}
