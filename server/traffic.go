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

package server

import (
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var requestCount int64

func TrafficStopMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUri := c.Request.RequestURI
		defer func() {
			if !strings.HasPrefix(requestUri, "/health") {
				atomic.AddInt64(&requestCount, -1)
			}
		}()

		if !strings.HasPrefix(requestUri, "/health") {
			atomic.AddInt64(&requestCount, 1)
		}

		c.Next()
	}

}

func requestsFinished() bool {
	return atomic.LoadInt64(&requestCount) == 0
}
