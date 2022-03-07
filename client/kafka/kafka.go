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
	"errors"
	"time"

	"github.com/Shopify/sarama"

	"github.com/NetEase-Media/ngo/adapter/log"
)

const (
	defaultVersion = "2.1.0"
)

var (
	consumerMap map[string]*Consumer
	producerMap map[string]*Producer
)

type Options struct {
	Name            string
	Addr            []string
	Version         string
	MaxOpenRequests int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Metadata        struct {
		Retries int
		Timeout time.Duration
	}
	Consumer struct {
		Group              string
		EnableAutoCommit   bool
		AutoCommitInterval time.Duration
		InitialOffset      int64
		SessionTimeout     time.Duration
		MinFetchBytes      int32
		DefaultFetchBytes  int32
		MaxFetchBytes      int32
		MaxFetchWait       time.Duration
		Retries            int
	}
	Producer struct {
		MaxMessageBytes  int
		Acks             sarama.RequiredAcks
		Timeout          time.Duration
		Retries          int
		MaxFlushBytes    int
		MaxFlushMessages int
		FlushFrequency   time.Duration
		Idempotent       bool
	}
}

func NewDefaultOptionsSlice(size int) []*Options {
	opts := make([]*Options, size)
	for i := 0; i < size; i++ {
		opts[i] = NewDefaultOptions()
	}
	return opts
}

func NewDefaultOptions() *Options {
	opt := &Options{}
	opt.Version = defaultVersion
	opt.MaxOpenRequests = 5
	opt.DialTimeout = time.Second * 30
	opt.ReadTimeout = time.Second * 30
	opt.WriteTimeout = time.Second * 30
	opt.Metadata.Retries = 3
	opt.Metadata.Timeout = time.Second * 60
	opt.Consumer.Group = ""
	opt.Consumer.EnableAutoCommit = true
	opt.Consumer.AutoCommitInterval = time.Second * 1
	opt.Consumer.InitialOffset = sarama.OffsetNewest
	opt.Consumer.SessionTimeout = time.Second * 10
	opt.Consumer.MinFetchBytes = 1
	opt.Consumer.DefaultFetchBytes = 1024 * 1024
	opt.Consumer.MaxFetchBytes = 0
	opt.Consumer.MaxFetchWait = time.Millisecond * 250
	opt.Consumer.Retries = 3
	opt.Producer.MaxMessageBytes = 1000000
	opt.Producer.Acks = sarama.WaitForLocal
	opt.Producer.Timeout = time.Second * 10
	opt.Producer.Retries = 3
	opt.Producer.MaxFlushBytes = 100 * 1024 * 1024
	opt.Producer.MaxFlushMessages = 0
	opt.Producer.FlushFrequency = time.Second * 1
	opt.Producer.Idempotent = false
	return opt
}

func Init(opts []*Options) error {
	if len(opts) == 0 {
		log.Info("empty kafka config, so skip init")
		return nil
	}
	sarama.Logger = New()

	for i := range opts {
		opt := opts[i]
		hasProducer := true // TODO default must have a producer
		hasConsumer := opt.Consumer.Group != ""
		if err := checkOptions(opt); err != nil {
			return err
		}
		if hasProducer {
			p, err := NewProducer(opt)
			if err != nil {
				return err
			}
			producerMap[opt.Name] = p
		}
		if hasConsumer {
			c, err := NewConsumer(opt)
			if err != nil {
				return err
			}
			consumerMap[opt.Name] = c
		}
	}
	return nil
}

func checkOptions(opt *Options) error {
	if opt.Version == "" {
		opt.Version = defaultVersion
	}

	if len(opt.Addr) == 0 {
		return errors.New("empty address")
	}
	return nil
}

func GetConsumer(name string) *Consumer {
	return consumerMap[name]
}

func GetProducer(name string) *Producer {
	return producerMap[name]
}

// StopAllConsumers 关闭所有kafka consumer
func StopAllConsumers() {
	for name, consumer := range consumerMap {
		consumer.Stop()
		log.Infof("Stop kafka consumer %s", name)
	}
}

// StopAllProducers 关闭所有kafka producer
func StopAllProducers() {
	for name, producer := range producerMap {
		producer.Close()
		log.Infof("Stop kafka producer %s", name)
	}
}

func StopAll() {
	StopAllConsumers()
	StopAllProducers()
}

func init() {
	consumerMap = make(map[string]*Consumer)
	producerMap = make(map[string]*Producer)

}
