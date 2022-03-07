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
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/xxl-job/xxl-job-executor-go"

	"github.com/NetEase-Media/ngo/adapter/protocol"
	"github.com/NetEase-Media/ngo/adapter/xxljob"
	"github.com/NetEase-Media/ngo/server"
)

// go run . -c ./config.yaml
func main() {
	s := server.Init()
	s.AddRoute(server.GET, "/hello", func(ctx *gin.Context) {
		ctx.JSON(protocol.JsonBody("hello"))
	})
	xxljob.RegTask("badjob", func(cxt context.Context, param *xxl.RunReq, logger *xxljob.XxlJobLogger) string {
		logger.Infof("badjob start...")
		logger.Infof("param :%+v", param)
		for i := 0; i < 1000; i++ {
			logger.Infof("what the hell:%d", i)
		}
		logger.Infof("badjob done...")
		return "ok"
	})
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("test response"))
	}))
	defer testServer.Close()

	s.Start()
}
