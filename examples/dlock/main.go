package main

import (
	"context"
	"sync"

	"github.com/NetEase-Media/ngo/pkg/dlock"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

// go run . -c ./app.yaml
func main() {
	ctx := context.Background()

	app := ngo.Init()
	app.PreStart = func() error {
		var n int
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				dlock.NewMutex("test", func() {
					n++
				}).WithTries(100).DoContext(ctx)
			}()
		}
		wg.Wait()
		log.Info("n is ", n)
		return nil
	}
	app.Start()
}

type Test struct {
	N1 int `json:"n1"`
	N2 int `json:"n2"`
}
