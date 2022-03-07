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
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/NetEase-Media/ngo/adapter/sentinel"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/adapter/protocol"
	"github.com/NetEase-Media/ngo/client/httplib"
	"github.com/NetEase-Media/ngo/server"
	"github.com/gin-gonic/gin"
)

// go run . -c ./config.yaml
func main() {
	s := server.Init()
	s.AddRoute(server.GET, "/hello", func(ctx *gin.Context) {
		ctx.JSON(protocol.JsonBody("hello"))
	})
	// log.GetLogger("haha").Infof("gungungun.....")
	// conf, err:= config.NewFromConfigFile("data.properties")
	// if err ==nil {
	// fmt.Println(conf)
	// }

	s.AddRoute(server.GET, "/error", func(ctx *gin.Context) {
		err := errors.New("db error")
		log.WithFields(
			"route", "error",
			"method", "get",
		).Errorf("send error %s", err.Error())
		ctx.JSON(protocol.ErrorJsonBody(protocol.DBError))
	})

	s.AddRoute(server.GET, "/panic", func(ctx *gin.Context) {
		panic(errors.New("panic test"))
	})

	s.AddRoute(server.POST, "/httplib", httpGetHandler)

	s.AddRoute(server.GET, "/cb", circuitbreakerHandler)

	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("test response"))
	}))
	defer testServer.Close()

	s.Start()
}

func httpGetHandler(ctx *gin.Context) {
	_, err := httplib.Get("http://www.baidu.com").Do(context.Background())
	if err != nil {
		log.Errorf("http get failed: %s", err.Error())
		ctx.JSON(protocol.ErrorJsonBody(protocol.SystemError))
	} else {
		ctx.JSON(protocol.JsonBody("Done!"))
	}
}

var (
	testServer *httptest.Server
)

func circuitbreakerHandler(ctx *gin.Context) {
	fakeError := errors.New("fake error")
	var res string
	statusCode, err := httplib.Get(testServer.URL).BindString(&res).CircuitBreaker("count", func() error {
		res = "circuitbreaker response"
		return fakeError
	}).Do(context.Background())

	var blockErr *sentinel.BlockError
	if errors.As(err, &blockErr) {
		log.Errorf("circuitbreaker block error %s", err.Error())
		log.Errorf("return my error? %v", errors.Is(err, fakeError))
		ctx.JSON(200, res) // circuitbreaker response
	} else {
		ctx.JSON(statusCode, res) // test response
	}

}
