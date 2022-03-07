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
	"sync"

	"github.com/NetEase-Media/ngo/dlock"

	"github.com/NetEase-Media/ngo/server"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	s := server.Init()
	s.AddRoute(server.GET, "/dlock", func(ctx *gin.Context) {
		var n1, n2 int
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				n1++
				dlock.NewMutex("test", func() {
					n2++
				}).WithTries(100).DoContext(ctx)
			}()
		}
		wg.Wait()
		ctx.JSON(200, &Test{n1, n2})
	})
	s.Start()
}

type Test struct {
	N1 int `json:"n1"`
	N2 int `json:"n2"`
}
