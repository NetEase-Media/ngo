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
	"encoding/json"

	"github.com/NetEase-Media/ngo/adapter/log"
)

// json序列化:将struct类型转换成string类型。 若不能序列化，返回""; 否则返回序列化后的string类型结果
// 为了输出方便，统一输出为string类型；若结果需要别的类型([]byte),请自己转换
func Marshal(str interface{}) string {
	ret, err := json.Marshal(str)
	if err != nil {
		log.Info("json不能序列化")
		return ""
	}
	return string(ret)
}

// json反序列化.将字符串转换成struct格式  若反序列化成功，则装入v中；不成功，则v中无数据
func Unmarshal(str string, v interface{}) {
	data := []byte(str)
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Info("json不能反序列化")
		return
	}
}
