package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/NetEase-Media/ngo/pkg/server/http/timeout"
	"github.com/gin-gonic/gin"
)

const defaultMsg = `{"code": -1, "msg":"http: Handler timeout"}`

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	s := http.Get()
	// create a handler that will last 1 seconds
	s.AddRoute(http.GET, "/short", timeoutFunc(short))

	// create a handler that will last 5 seconds
	s.AddRoute(http.GET, "/long", timeoutFunc(long))

	// create a handler that will last 5 seconds but can be canceled.
	s.AddRoute(http.GET, "/long2", timeoutFunc(long2))

	s.AddRoute(http.GET, "/boundary", timeoutFunc(boundary))

	s.AddRoute(http.GET, "/panic", timeoutFunc(panic))
	app.Start()
}

func timeoutFunc(f gin.HandlerFunc) gin.HandlerFunc {
	return timeout.Timeout(
		timeout.WithTimeout(2*time.Second),
		timeout.WithHandler(f),
		timeout.WithErrorHttpCode(http.StatusServiceUnavailable), // optional
		timeout.WithDefaultMsg(defaultMsg),                       // optional
		timeout.WithErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": -1})
		}),
		/*
			timeout.WithCallBack(func(r *http.Request) {
				fmt.Println("timeout happen, url:", r.URL.String())
			}) // optional
		*/
	)
}

func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "short"})
}

func long(c *gin.Context) {
	time.Sleep(3 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "long"})
}

func boundary(c *gin.Context) {
	time.Sleep(2 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "boundary"})
}

func long2(c *gin.Context) {
	if doSomething(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{"hello": "long2"})
	}
}

func panic(c *gin.Context) {
	time.Sleep(1 * time.Second)
	x := 0
	fmt.Println(100 / x)
}

// A cancelCtx can be canceled.
// When canceled, it also cancels any children that implement canceler.
func doSomething(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		fmt.Println("doSomething is canceled.")
		return false
	case <-time.After(5 * time.Second):
		fmt.Println("doSomething is done.")
		return true
	}
}
