package main

import (
	"strconv"

	"github.com/NetEase-Media/ngo/pkg/client/kafka"
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/ngo"
)

const topic = "test"

var (
	consumer    *kafka.Consumer
	producer    *kafka.Producer
	messageChan = make(chan kafka.ConsumerMessage, 1000)
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	app.PreStart = func() error {
		consumer = kafka.GetConsumer("k1")
		//consumer.AddListener(topic, &listener{})
		consumer.AddBatchListener(topic, &batchListener{})
		consumer.Start()

		producer = kafka.GetProducer("k1")
		for i := 0; i < 98; i++ {
			producer.Send(topic, "hello world!"+strconv.Itoa(i), nil)
		}
		go func() {
			for {
				r := <-messageChan
				log.Info(r.Value, " ", r.Partition, " ", r.Offset)
			}
		}()
		return nil
	}
	app.Start()
}

type listener struct {
	kafka.Listener
}

func (l *listener) Listen(message kafka.ConsumerMessage, ack *kafka.Acknowledgment) {
	messageChan <- message
	ack.Acknowledge()
}

type batchListener struct {
	kafka.BatchListener
}

func (l *batchListener) Listen(messages []kafka.ConsumerMessage, ack *kafka.Acknowledgment) {
	log.Info("batch size is ", len(messages))
	for _, msg := range messages {
		messageChan <- msg
	}
	ack.Acknowledge()
}

func (l *batchListener) BatchCount() int {
	return 10
}
