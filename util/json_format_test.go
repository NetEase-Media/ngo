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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 分别对json进行序列化和反序列测试
type Stu struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	High string `json:"high"`
	Sex  string `json:"sex"`
}

//序列化成功
func TestMarshal(t *testing.T) {
	stu := Stu{
		Name: "陈奕迅",
		Age:  47,
		High: "173cm",
		Sex:  "男",
	}
	expect := "{\"name\":\"陈奕迅\",\"age\":47,\"high\":\"173cm\",\"sex\":\"男\"}"
	actual := Marshal(stu)
	assert.Equal(t, expect, actual, "json序列化不符合预期")
}

//序列化失败,err!=nil
func TestMarshalErr(t *testing.T) {
	c := make(chan int, 1)
	cur := Marshal(c)
	assert.Equal(t, "", cur, "")
}

//反序列化成功
func TestUnmarsha(t *testing.T) {
	actual := Stu{}
	expected := Stu{
		Name: "陈奕迅",
		Age:  47,
		High: "173cm",
		Sex:  "男",
	}
	data := "{\"name\":\"陈奕迅\",\"age\":47,\"high\":\"173cm\",\"sex\":\"男\"}"
	Unmarshal(data, &actual)
	fmt.Println(actual)
	assert.Equal(t, expected, actual, "")
}

//反序列化失败，err!=nil
func TestUnmarshaErr(t *testing.T) {
	ss := Stu{}
	stu := Stu{}
	data := "{\"name\":\"陈奕迅\",\"age\":47,\"high\":\"173cm\",\"sex\":\"男\""
	Unmarshal(data, &stu)
	assert.Equal(t, ss, stu, "")
}
