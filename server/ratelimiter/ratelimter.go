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

package ratelimiter

import (
	"net/http"

	"github.com/NetEase-Media/ngo/adapter/sentinel"
	"github.com/gin-gonic/gin"
)

func NewDefaultOptions() *RatelimiterOptions {
	return &RatelimiterOptions{
		DefaultMsg:    `{"code": -1, "msg":"http: Handler limit"}`,
		ErrorHttpCode: http.StatusTooManyRequests,
	}
}

func RateLimiter(opts ...Option) gin.HandlerFunc {
	o := NewDefaultOptions()

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(o)
	}

	return func(c *gin.Context) {
		e, b := sentinel.Entry(o.Resource)
		if b != nil {
			if o.ErrorHandler != nil {
				o.ErrorHandler(c)
				c.Abort()
			} else {
				c.AbortWithStatusJSON(o.ErrorHttpCode, o.DefaultMsg)
				return
			}
		} else {
			e.Exit()
			c.Next()
		}
	}
}
