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

package zookeeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 先使用kafka的zk测试
var testZkClientAddr = "kafka"

func TestInit(t *testing.T) {
	var opts []Options

	o := Options{
		Name:           "m",
		Addr:           []string{testZkClientAddr},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	err := Init(opts)
	assert.Equal(t, nil, err, "encounter error.")
}

//func TestInitException1(t *testing.T) {
//	var opts []Options
//	o := Options{
//		Name:           "m",
//		Addr:           []string{ADDR},
//		SessionTimeout: time.Second * 5,
//	}
//	opts = append(opts, o)
//	err := Init(opts)
//	assert.Equal(t, nil, err, "encounter error.")
//	assert.Panics(t, func() {
//		Init(opts)
//	})
//}
//func TestInitException2(t *testing.T) {
//	var opts []Options
//	err := Init(opts)
//	assert.Equal(t, nil, err, "encounter error.")
//}

func TestNewClientFromOptionException1(t *testing.T) {
	var opts []Options
	o := Options{
		Name:           "",
		Addr:           []string{testZkClientAddr},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	//name=""
	_, err := NewClientFromOption(opts)
	assert.Error(t, err)
}

func TestNewClientFromOptionException2(t *testing.T) {
	var opts []Options
	o := Options{
		Name:           "m",
		Addr:           []string{},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	_, err := NewClientFromOption(opts)
	assert.Error(t, err)
}

func TestNewClientFromOptionException3(t *testing.T) {
	var opts []Options
	o := Options{
		Name:           "m",
		Addr:           []string{"a"},
		SessionTimeout: time.Second * 1,
	}
	opts = append(opts, o)
	_, err := NewClientFromOption(opts)
	assert.Error(t, err)
}
func TestNewClientFromOptionException4(t *testing.T) {
	var opts []Options
	o := Options{
		Name:           "m",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	o1 := Options{
		Name:           "m",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o1)
	_, err := NewClientFromOption(opts)
	assert.Error(t, err)
}
