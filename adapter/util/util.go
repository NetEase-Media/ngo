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
	"reflect"
)

// CheckError 提供简介的error判断，如果err != nil则panic
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// MaxInt64 返回数组中最大的
func MaxInt64(i int64, nums ...int64) int64 {
	max := i
	for _, num := range nums {
		if max < num {
			max = num
		}
	}
	return max
}

// MinInt64 返回数组中最小的
func MinInt64(i int64, nums ...int64) int64 {
	min := i
	for _, num := range nums {
		if min > num {
			min = num
		}
	}
	return min
}

func TypeName(data interface{}) string {
	if data == nil {
		return ""
	}
	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr { // 解引用嵌套指针
		t = t.Elem()
	}

	return t.String()
}
