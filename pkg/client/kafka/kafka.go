package kafka

import "github.com/Shopify/sarama"

const (
	defaultVersion = "2.1.0"
)

func init() {
	sarama.Logger = NewLogger()
}

func New(opt *Options) (*Kafka, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	k := &Kafka{
		Opt: opt,
	}
	hasProducer := true // TODO default must have a producer
	hasConsumer := opt.Consumer.Group != ""
	if hasProducer {
		p, err := NewProducer(opt)
		if err != nil {
			return nil, err
		}
		k.Producer = p
	}
	if hasConsumer {
		co, err := NewConsumer(opt)
		if err != nil {
			return nil, err
		}
		k.Consumer = co
	}
	return k, nil
}

type Kafka struct {
	Opt      *Options
	Consumer *Consumer
	Producer *Producer
}
