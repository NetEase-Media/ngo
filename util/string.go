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

import "strings"

// Symbol 分隔符
type Symbol string

// 标记
const (
	Point             Symbol = "."
	Comma             Symbol = ","
	Colon             Symbol = ":"
	Empty             Symbol = ""
	Blank             Symbol = " "
	Underline         Symbol = "_"
	DoubleVerticalBar Symbol = "||"
	VerticalBar       Symbol = "|"
	QuestionMark      Symbol = "?"
	Percent           Symbol = "%"
	Enter             Symbol = "\n"
)

// Join 通过分割符号将数组转换成字符串
func Join(list []string, symbol Symbol) string {
	if len(list) == 0 {
		return string(Empty)
	}
	var l []string
	for _, v := range list {
		// if v == "" { // 过滤掉为空的数据
		// 	continue  // 先去掉兼容java工程
		// }
		l = append(l, v)
	}
	r := strings.Join(l, string(symbol))
	return strings.Trim(r, string(Blank))
}

// Split 将字符串根据分割符号转换成数组
func Split(s string, symbol Symbol) []string {
	var r []string = make([]string, 0)

	if s == "" {
		return r
	}
	l := strings.Split(s, string(symbol))
	if len(l) == 0 {
		return r
	}
	for _, v := range l {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		r = append(r, v)
	}
	return r
}

// SplitNoRepeat 无重复的数组
func SplitNoRepeat(s string, symbol Symbol) []string {
	var r []string = make([]string, 0)

	if s == "" {
		return r
	}

	l := strings.Split(s, string(symbol))
	if len(l) == 0 {
		return r
	}

	var checker map[string]int = make(map[string]int, len(l))
	for _, v := range l {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		_, ok := checker[v]
		if ok {
			continue
		}
		checker[v] = 0 // 设置检查位
		r = append(r, v)
	}
	return r
}
