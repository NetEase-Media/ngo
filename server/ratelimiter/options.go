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
	"github.com/gin-gonic/gin"
)

type Option func(*RatelimiterOptions)

type RatelimiterOptions struct {
	Resource      string
	DefaultMsg    string
	ErrorHttpCode int
	ErrorHandler  gin.HandlerFunc
}

// Optional parameters
func WithResource(resource string) Option {
	return func(o *RatelimiterOptions) {
		o.Resource = resource
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(o *RatelimiterOptions) {
		o.ErrorHttpCode = code
	}
}

// Optional parameters
func WithDefaultMsg(s string) Option {
	return func(o *RatelimiterOptions) {
		o.DefaultMsg = s
	}
}

// Optional parameters
func WithErrorHandler(f gin.HandlerFunc) Option {
	return func(o *RatelimiterOptions) {
		o.ErrorHandler = f
	}
}
