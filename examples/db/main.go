package main

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/client/db"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		db := db.WithContext(context.Background(), db.GetClient("ngo"))
		var t test
		db.Raw("select blackword from blacklist t where bid = ?", 1).Find(&t)
		log.Info(t.Blackword)

		db.Exec("update blacklist set blackword = ? where bid = ?", "mmm", 21)
		return nil
	}
	app.Start()
}

type test struct {
	Blackword string `json:"blackword"`
}
