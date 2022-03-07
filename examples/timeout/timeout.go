// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/NetEase-Media/ngo/server"

	"github.com/NetEase-Media/ngo/server/timeout"
	"github.com/gin-gonic/gin"
)

const defaultMsg = `{"code": -1, "msg":"http: Handler timeout"}`

// go run . -c ./app.yaml
func main() {

	s := server.Init()
	// add timeout middleware with 2 second duration

	// create a handler that will last 1 seconds
	s.AddRoute(server.GET, "/short", timeoutFunc(short))

	// create a handler that will last 5 seconds
	s.AddRoute(server.GET, "/long", timeoutFunc(long))

	// create a handler that will last 5 seconds but can be canceled.
	s.AddRoute(server.GET, "/long2", timeoutFunc(long2))

	s.AddRoute(server.GET, "/boundary", timeoutFunc(boundary))

	s.AddRoute(server.GET, "/panic", timeoutFunc(panic))

	s.Start()
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
