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
	"encoding/xml"

	"github.com/NetEase-Media/ngo/adapter/log"
)

// xml序列化：将struct类型序列化成xml类型。 若成功，则输出xml的字符串；否则输出""
// 为了输出方便，统一输出为string类型；若结果需要别的类型([]byte),请自己转换
func MarshalIndent(v interface{}, prefix, indent string) string {
	data, err := xml.MarshalIndent(v, prefix, indent)
	if err != nil {
		log.Info("xml不能序列化")
		return ""
	}
	return string(data)
}

// xml的反序列化： 将xml字符串转换成序列化。 若反序列化成功，则将数据放到v中；否则v中无数据
func UnmarshalIndent(str string, v interface{}) {
	data := []byte(str)
	err := xml.Unmarshal(data, v)
	if err != nil {
		log.Info("xml不能反序列化")
		return
	}
}
