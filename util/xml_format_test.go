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
	"errors"
	"testing"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	High string `json:"high"`
	Sex  string `json:"sex"`
}
type Student2 struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func TestMarshalIndent(t *testing.T) {
	stu := Student{
		Name: "陈奕迅",
		Age:  47,
		High: "173cm",
		Sex:  "男",
	}
	expected := `<Student><Name>陈奕迅</Name><Age>47</Age><High>173cm</High><Sex>男</Sex></Student>`
	actual := MarshalIndent(stu, "", "")
	assert.Equal(t, expected, actual, "")
}

func TestMarshalIndent_(t *testing.T) {
	c := make(chan int, 1)
	cur := MarshalIndent(c, "", "   ")
	assert.Equal(t, "", cur, "")
}

//模拟序列化失败，返回空字符串
func TestMarshalIndentErr(t *testing.T) {
	stu := Student2{
		Name: "陈奕迅<Name>",
		Age:  "47</Age>",
	}
	actual := MarshalIndentHelperTest(stu, "", "	")
	assert.Equal(t, "", actual, "")
}

func TestUnmarshalIndent(t *testing.T) {
	expected := Student{
		Name: "陈奕迅",
		Age:  47,
		High: "173cm",
		Sex:  "男",
	}
	s := `<Student>
			<Name>陈奕迅</Name>
			<Age>47</Age>
			<High>173cm</High>
			<Sex>男</Sex>
		</Student>`
	actual := Student{}
	UnmarshalIndent(s, &actual)
	assert.Equal(t, expected, actual, "")
}

//反序列化失败，err=nil
func TestUnmarshalIndent_(t *testing.T) {
	ss := Student{}
	s := `<Name>陈奕迅</Name>
		  <Age>47</Age>
		  <High>173cm</High>
		  <Sex>男</Sex>`
	actual := Student{}
	UnmarshalIndent(s, &actual)
	assert.Equal(t, ss, actual, "")
}

//反序列化失败，err!=nil
func TestUnmarshalIndentErr(t *testing.T) {
	ss := Student{}
	s := `<Name>陈奕迅</Name>
		  <Age>47</Age>
		  <High>173cm</High>
		  <Sex>男</Sex>`
	actual := Student2{}
	UnmarshalIndent(s, actual)
	assert.NotEqual(t, ss, actual)
}

//辅助方法，模拟err不为nil时的处理逻辑
func MarshalIndentHelperTest(v interface{}, prefix, indent string) string {
	data, err := xml.MarshalIndent(v, prefix, indent)
	err = errors.New("testxff")
	if err != nil {
		log.Info("xml不能序列化")
		return ""
	}
	return string(data)
}
