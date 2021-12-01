//+build !race

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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTrafficStop(t *testing.T) {
	r := gin.New()
	r.Use(TrafficStopMiddleware())

	r.GET("/", func(context *gin.Context) {
		time.Sleep(1000 * time.Millisecond)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	go r.ServeHTTP(w, req)

	time.Sleep(500 * time.Millisecond)
	assert.False(t, requestsFinished())
	time.Sleep(800 * time.Millisecond)
	assert.True(t, requestsFinished())
}
