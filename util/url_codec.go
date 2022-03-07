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
	"net/url"
	"strings"
)

// Encode encode
func Encode(query string) string {
	return url.QueryEscape(query)
}

// Decode decode
func Decode(query string) string {
	rs, err := url.QueryUnescape(query)
	if err != nil { // 解析出现问题，直接不处理
		return query
	}
	return rs
}

//EncodeEscape 处理<p>a=10&b=100&c=</p>
// 但是其中的'&' and '=' 不做处理
// 原生的方式自己可以通过<p>url.Values</p>自行处理
func EncodeEscape(query string) string {
	val := parse(query)
	return val.Encode()
}

func parse(query string) url.Values {
	rs := make(url.Values)
	if query == "" {
		return rs
	}

	kvs := strings.Split(query, "&")
	for _, kv := range kvs {
		vs := strings.Split(kv, "=")
		if len(vs) > 1 {
			rs.Add(vs[0], vs[1])
		} else if len(vs) == 1 {
			rs.Add(vs[0], "")
		}
	}
	return rs
}
