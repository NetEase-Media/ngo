package main

import (
	"github.com/NetEase-Media/ngo/pkg/client/zookeeper"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

var client *zookeeper.ZookeeperProxy

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()

	app.PreStart = func() error {
		client = zookeeper.GetClient("zookeeper01")
		path, err := client.CreateNode("/testt", zookeeper.EPHEMERAL_SEQUENTIAL, "test")
		if err != nil {
			log.Error(err)
		} else {
			log.Info(path)
		}
		return nil
	}
	app.Start()
}
