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

package main

import (
	"net/http"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/gin-gonic/gin"

	"github.com/NetEase-Media/ngo/client/kafka"
	"github.com/NetEase-Media/ngo/server"
)

const topic = "test"

var (
	consumer    *kafka.Consumer
	producer    *kafka.Producer
	messageChan chan kafka.ConsumerMessage
)

// go run . -c ./app.yaml
func main() {
	s := server.Init()

	s.PreStart = func() error {
		producer = kafka.GetProducer("k1")

		messageChan = make(chan kafka.ConsumerMessage, 1000)
		consumer = kafka.GetConsumer("k1")
		consumer.AddListener(topic, &listener{})
		consumer.Start()
		return nil
	}

	s.AddRoute(server.POST, "/put-async", func(ctx *gin.Context) {
		v := ctx.PostForm("value")
		if v == "" {
			ctx.JSON(http.StatusBadRequest, "")
			return
		}
		//
		producer.Send(topic, v, func(metadata *kafka.RecordMetadata, err error) {
			if err != nil {
				log.Errorf("send error: %v", err)
			}
		})
		ctx.JSON(200, v)
	})

	s.AddRoute(server.POST, "/put-sync", func(ctx *gin.Context) {
		v := ctx.PostForm("value")
		if v == "" {
			ctx.JSON(http.StatusBadRequest, "")
			return
		}
		//
		err := producer.SyncSend(topic, v)
		if err != nil {
			log.Errorf("send error: %v", err)
		}
		ctx.JSON(200, v)
	})

	s.AddRoute(server.GET, "/get", func(ctx *gin.Context) {
		select {
		case m := <-messageChan:
			ctx.JSON(200, m.Value)
		default:
			ctx.JSON(http.StatusInternalServerError, "no data")
		}
	})

	s.Start()
}

type listener struct {
	kafka.Listener
}

func (l *listener) Listen(message kafka.ConsumerMessage, ack *kafka.Acknowledgment) {
	messageChan <- message
	ack.Acknowledge()
}
