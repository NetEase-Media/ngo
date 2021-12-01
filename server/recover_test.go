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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/NetEase-Media/ngo/adapter/log"
)

func doPanic(ctx *gin.Context) {
	log.Panic("panic...")
}

func TestRecover(t *testing.T) {
	r := gin.New()
	r.Use(ServerRecover())

	r.GET("/", doPanic)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestRecover2(t *testing.T) {
	r := gin.New()
	r.Use(OutermostRecover())

	r.GET("/", doPanic)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", w.Body.String())
}
