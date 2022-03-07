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
	"net"
	"net/http"
	"strings"
)

// 获取用户IP
func GetRequestIp(r *http.Request) string {

	// 用户经过代理时会配置报头X-Real-IP为用户的真实IP
	ip := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}

	// 用户经过代理时会在X-Forwarded-For后面添加上一级设备的IP
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	// 通过RemoteAddr获取与服务端直接相连的设备IP
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return ip
	}

	return ""

}
