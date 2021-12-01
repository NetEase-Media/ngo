//+build !race

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

package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	NAME = "ngo-test-init"
)

//func TestInit_AddrNil(t *testing.T) {
//	opts := NewDefaultOptionsSlice(1)
//	opts[0].Addr = nil
//	err := Init(opts)
//	assert.Error(t, err)
//}

//func TestInit_optsNil(t *testing.T) {
//	var opts []*Options
//	err := Init(opts)
//	assert.Equal(t, nil, err)
//}

func TestInit_InitProcess(t *testing.T) {
	opts := NewDefaultOptionsSlice(1)
	opts[0].Name = NAME
	opts[0].Addr = []string{KAFKAADDR}
	opts[0].Version = KAFKAVERSION
	opts[0].Consumer.Group = "ngo"
	err := Init(opts)
	assert.Equal(t, nil, err)
	consumer := GetConsumer(NAME)
	assert.Equal(t, "ngo", consumer.opt.Consumer.Group)
	producer := GetProducer(NAME)
	assert.Equal(t, time.Second*10, producer.opt.Producer.Timeout)
}

//func TestInit_InitProcess_groupEmpty(t *testing.T) {
//	opts := NewDefaultOptionsSlice(1)
//	opts[0].Name = NAME
//	opts[0].Addr = []string{KAFKAADDR}
//	opts[0].Version = KAFKAVERSION
//	err := Init(opts)
//	assert.Equal(t, nil, err)
//	consumer := GetConsumer(NAME)
//	assert.Nil(t, consumer)
//	producer := GetProducer(NAME)
//	assert.Equal(t, time.Second*10, producer.opt.Producer.Timeout)
//}

//func TestInit_CloseConsumer(t *testing.T) {
//	opts := NewDefaultOptionsSlice(1)
//	opts[0].Name = NAME
//	opts[0].Addr = []string{KAFKAADDR}
//	opts[0].Version = KAFKAVERSION
//	opts[0].Consumer.Group = "ngo"
//	err := Init(opts)
//	assert.NoError(t, err)
//	producer := GetProducer(NAME)
//	go func() {
//		i := 0
//		for {
//			producer.Send("ngo-test-kafka-init", "init-test-ahah-"+strconv.Itoa(i), nil)
//			i++
//		}
//	}()
//	consumer := GetConsumer(NAME)
//	assert.Equal(t, "ngo", consumer.opt.Consumer.Group)
//	consumer.AddListener("ngo-test-kafka-init", &listener{func(message ConsumerMessage, ack *Acknowledgment) {
//		assert.Contains(t, message.Value, "init-test-ahah-")
//		//log.Info(message.Value)
//	}})
//	consumer.Start()
//	time.Sleep(time.Second * 3)
//	StopAllConsumers()
//	time.Sleep(time.Second * 5)
//}
