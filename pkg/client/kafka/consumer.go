package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/metrics"
	"github.com/Shopify/sarama"
)

type ConsumerMessage struct {
	Topic     string
	Key       string
	Value     string
	Partition int32
	Offset    int64
}

type Listener interface {
	Listen(ConsumerMessage, *Acknowledgment)
}

type BatchListener interface {
	Listen([]ConsumerMessage, *Acknowledgment)
	BatchCount() int
}

// Consumer 是一个group的消费者
type Consumer struct {
	client         sarama.ConsumerGroup
	logger         log.Logger
	opt            Options
	ctx            context.Context
	cancel         func()
	runChan        chan struct{}
	listeners      map[string]Listener
	batchListeners map[string]BatchListener
}

func (c *Consumer) Options() Options {
	return c.opt
}

func (c *Consumer) AddListener(topic string, listener Listener) {
	if len(topic) == 0 {
		panic("topic must not be empty")
	}
	if listener == nil {
		panic("listener must not be nil")
	}
	c.listeners[topic] = listener
}

func (c *Consumer) AddBatchListener(topic string, listener BatchListener) {
	if len(topic) == 0 {
		panic("topic must not be empty")
	}
	if listener == nil {
		panic("listener must not be nil")
	}
	c.batchListeners[topic] = listener
}

// Start 启动后台消费任务
func (c *Consumer) Start() {
	if len(c.listeners) == 0 && len(c.batchListeners) == 0 {
		panic("empty topic listener")
	}

	// 当前不允许多个后台消费任务
	if c.ctx != nil {
		panic("duplicated start")
	}

	h := &consumerHandler{
		consumer: c,
		ready:    make(chan struct{}),
		logger:   c.logger,
		opt:      &c.opt,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.runChan = make(chan struct{})
	tmap := make(map[string]struct{}, len(c.listeners)+len(c.batchListeners))
	for k := range c.listeners {
		tmap[k] = struct{}{}
	}
	for k := range c.batchListeners {
		tmap[k] = struct{}{}
	}

	topics := make([]string, 0, len(tmap))
	for k := range tmap {
		topics = append(topics, k)
	}

	go func() {
		defer close(c.runChan)
		for {
			// 当服务的rebalance后会返回
			if err := c.client.Consume(c.ctx, topics, h); err != nil {
				log.Errorf("kafka consume failed: %s", err.Error())
				time.Sleep(time.Millisecond * 200) // 睡眠防止异常之后死循环占满CPU
			}

			if c.ctx.Err() != nil {
				return
			}

			select {
			case <-h.ready:
				h.ready = make(chan struct{})
			default:
			}
		}
	}()
	<-h.ready
	c.logger.Info("consumer up and running")
}

// Stop 停止后台消费任务
func (c *Consumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
		<-c.runChan
	}
	return c.client.Close()
}

func NewConsumer(opt *Options) (*Consumer, error) {
	config, err := newConsumerConfig(opt)
	if err != nil {
		return nil, err
	}
	c, err := sarama.NewConsumerGroup(opt.Addr, opt.Consumer.Group, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		client: c,
		logger: log.WithFields(
			"kafka", opt.Name,
			"group", opt.Consumer.Group,
		),
		opt:            *opt,
		listeners:      make(map[string]Listener, 8),
		batchListeners: make(map[string]BatchListener, 8),
	}, nil
}

func newConsumerConfig(opt *Options) (*sarama.Config, error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(opt.Version)
	if err != nil {
		return nil, err
	}
	config.Version = version

	config.Metadata.RefreshFrequency = time.Second * 10

	config.Net.MaxOpenRequests = opt.MaxOpenRequests
	config.Net.DialTimeout = opt.DialTimeout
	config.Net.ReadTimeout = opt.ReadTimeout
	config.Net.WriteTimeout = opt.WriteTimeout

	config.Net.SASL.Enable = opt.SASL.Enable
	config.Net.SASL.Mechanism = opt.SASL.Mechanism
	config.Net.SASL.User = opt.SASL.User
	config.Net.SASL.Password = opt.SASL.Password
	config.Net.SASL.Handshake = opt.SASL.Handshake

	config.Metadata.Retry.Max = opt.Metadata.Retries
	config.Metadata.Timeout = opt.Metadata.Timeout

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = opt.Consumer.EnableAutoCommit
	config.Consumer.Offsets.AutoCommit.Interval = opt.Consumer.AutoCommitInterval
	config.Consumer.Offsets.Initial = opt.Consumer.InitialOffset
	config.Consumer.Offsets.Retry.Max = opt.Consumer.Retries
	config.Consumer.Group.Session.Timeout = opt.Consumer.SessionTimeout
	config.Consumer.Fetch.Min = opt.Consumer.MinFetchBytes
	config.Consumer.Fetch.Default = opt.Consumer.DefaultFetchBytes
	config.Consumer.Fetch.Max = opt.Consumer.MaxFetchBytes
	config.Consumer.MaxWaitTime = opt.Consumer.MaxFetchWait
	return config, nil
}

// consumerHandler 用来运行消费者后台任务
type consumerHandler struct {
	consumer *Consumer
	ready    chan struct{}
	logger   log.Logger
	opt      *Options
}

// Setup 在启动前执行
func (ch *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(ch.ready)
	return nil
}

// Cleanup 在结束后执行
func (ch *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 在循环中消费message
func (ch *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	//
	if listener := ch.consumer.listeners[topic]; listener != nil {
		for message := range claim.Messages() {
			ch.logger.Tracef("Message claimed: value = %s, timestamp = %v, topic = %s",
				string(message.Value), message.Timestamp, message.Topic)
			ch.listen(listener, session, message)
		}
	} else if batchListener := ch.consumer.batchListeners[topic]; batchListener != nil {
		count := batchListener.BatchCount()
		if count < 1 {
			count = 1
		}
		msgArr := make([]*sarama.ConsumerMessage, 0, count)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case message := <-claim.Messages():
				if message == nil {
					log.Info("channel is closed")
					return nil
				}
				msgArr = append(msgArr, message)
				if len(msgArr) >= count {
					ch.batchListen(batchListener, session, msgArr)
					msgArr = make([]*sarama.ConsumerMessage, 0, count)
				}
			case <-ticker.C:
				if len(msgArr) > 0 {
					ch.batchListen(batchListener, session, msgArr)
					msgArr = make([]*sarama.ConsumerMessage, 0, count)
				}
			}
		}
	}
	return nil
}

func (ch *consumerHandler) listen(listener Listener, session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	msg := ConsumerMessage{
		Topic:     message.Topic,
		Key:       string(message.Key),
		Value:     string(message.Value),
		Partition: message.Partition,
		Offset:    message.Offset,
	}
	ack := &Acknowledgment{
		ch:      ch,
		session: session,
		message: message,
	}
	begin := time.Now()
	defer func() {
		var err error
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("unexpected panic value: %#v", r)
		}
		if err != nil {
			json, _ := json.Marshal(&msg)
			log.Errorf("consumer handle error: %v, message: %s", err, json)
		}
		ch.collect(message.Topic, message.Partition, len(message.Value), time.Since(begin), err)
	}()

	listener.Listen(msg, ack)
	// if auto commit, mark message
	if ch.consumer.opt.Consumer.EnableAutoCommit {
		session.MarkMessage(message, "")
	}
}

func (ch *consumerHandler) batchListen(listener BatchListener, session sarama.ConsumerGroupSession, msgArr []*sarama.ConsumerMessage) {
	begin := time.Now()
	topic := msgArr[0].Topic
	partition := msgArr[0].Partition
	var msgBytes int
	defer func() {
		var err error
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("unexpected panic value: %#v", r)
		}
		if err != nil {
			log.Errorf("batch consumer handle error: %v", err)
		}
		ch.collect(topic, partition, msgBytes, time.Since(begin), err)
	}()
	msgs := make([]ConsumerMessage, 0, len(msgArr))
	for _, message := range msgArr {
		msgs = append(msgs, ConsumerMessage{
			Topic:     message.Topic,
			Key:       string(message.Key),
			Value:     string(message.Value),
			Partition: message.Partition,
			Offset:    message.Offset,
		})
		msgBytes += len(message.Value)
	}
	ack := &Acknowledgment{
		ch:      ch,
		session: session,
		message: msgArr[len(msgArr)-1],
	}
	listener.Listen(msgs, ack)

	// if auto commit, mark message
	if ch.consumer.opt.Consumer.EnableAutoCommit {
		session.MarkMessage(msgArr[len(msgArr)-1], "")
	}
}

// collect 生成监控数据发送到收集器
func (ch *consumerHandler) collect(topic string, partition int32, msgBytes int, cost time.Duration, err error) {
	if !metrics.IsMetricsEnabled() {
		return
	}
}

type Acknowledgment struct {
	ch      *consumerHandler
	session sarama.ConsumerGroupSession
	message *sarama.ConsumerMessage
}

func (a *Acknowledgment) Acknowledge() {
	if !a.ch.consumer.opt.Consumer.EnableAutoCommit {
		a.session.MarkMessage(a.message, "")
		a.session.Commit()
	}
}
