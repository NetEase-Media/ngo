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
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NetEase-Media/ngo/adapter/job"
	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/server"
)

// go run . -c ./config.yaml
func main() {
	s := server.Init()

	// 两个接口用来模拟调度平台，开发时无需实现
	s.AddRoute(server.GET, "/api/v1/:namespace/podshared/getSharedNum", getSharedNumHandler)
	s.AddRoute(server.POST, "/api/v1/:namespace/podshared/report", circuitbreakerHandler)

	go s.Start()

	job.Run(func(args *job.Args) (string, error) {
		var b strings.Builder
		for i := 0; i <= args.SharedNum; i++ {
			b.Write([]byte(`
    /\_____/\
   /  o   o  \
  ( ==  ^  == )
   )         (
  (           )
 ( (  )   (  ) )
(__(__)___(__)__)

`))
		}
		return b.String(), nil
	})
}

func getSharedNumHandler(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	podName := ctx.Query("podName")
	log.Infof("receive getSharedNum request namespace %s podName %s", namespace, podName)

	res := &job.GetShareNumResponse{
		Status:  0,
		Message: "",
		Data: job.GetShareNumData{
			SharedNum: 2,
		},
	}
	ctx.JSON(200, res)
}

func circuitbreakerHandler(ctx *gin.Context) {
	namespace := ctx.Param("namespace")

	var req job.ReportRequest
	ctx.BindJSON(&req)
	log.Infof("receive report namespace %s podName %s result %s", namespace, req.PodName, req.Result)

	ctx.JSON(200, &job.ReportResponse{Status: 200})
}
