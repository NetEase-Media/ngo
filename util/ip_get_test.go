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

package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestGetRequestIp(t *testing.T) {
	r := gin.Context{Request: &http.Request{Header: http.Header{}}}
	//ip为空
	ip_ := GetRequestIp(r.Request)
	assert.Equal(t, "", ip_)
	r.Request.Header.Set("X-Real-IP", " 10.10.10.10 ")
	r.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	r.Request.RemoteAddr = "  40.40.40.40:42123 "

	if ip := GetRequestIp(r.Request); ip != "10.10.10.10" {
		t.Errorf("actual: 10.10.10.10, expected:%s", ip)
	}

	r.Request.Header.Del("X-Real-IP")
	if ip := GetRequestIp(r.Request); ip != "20.20.20.20" {
		t.Errorf("actual: 20.20.20.20, expected:%s", ip)
	}

	r.Request.Header.Set("X-Real-IP", "30.30.30.30  ")
	if ip := GetRequestIp(r.Request); ip != "30.30.30.30" {
		t.Errorf("actual: 30.30.30.30, expected:%s", ip)
	}

	r.Request.Header.Del("X-Forwarded-For")
	r.Request.Header.Del("X-Real-IP")
	if ip := GetRequestIp(r.Request); ip != "220.220.220.220" {
		t.Errorf("actual: 220.220.220.220, expected:%s", ip)
	}

	r.Request.RemoteAddr = "50.50.50.50"
	if ip := GetRequestIp(r.Request); ip != "" {
		t.Errorf("ip: 50.50.50.50")
	}
}
