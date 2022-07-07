# [Ngo](https://github.com/NetEase-Media/ngo)

---
## kafka client
### 模块用途
通过kafka的生产、消费功能，Kafka客户端在 [sarama](https://github.com/Shopify/sarama) 基础上实现，由于原始接口比较复杂，业务需求一般用不上，Ngo中对其进行了较多的封装。在配置文件中增加kafka段，Ngo即会自动按配置生成生产者和消费者。
### 使用说明
#### 配置
参数配置详见[kafka options](config.md#kafka-配置-kafkaoptions)
#### 生产者
##### 获取实例
```go
// 配置文件方式
p := kafka.GetProducepr("p")

// 手动实例化
opts := kafka.NewDefaultOptions()
opts.Addr = []string{"kafka:9092"}
p, err := kafka.NewProducer(opts)
```
##### 发送消息
```go
// 异步发送简单消息
p.Send("topic1", "message1", func(metadata *kafka.RecordMetadata, err error){
	if (err != nil) {
		log.Error("send error", err)
	}   
})

// 异步发送消息，指定key，相同key在同一分区
p.SendMessage(kafka.ProducerMessage{Topic: "topic1", Key: "key1", Value: "value1"}, func(err error){})

// 同步发送简单消息
err := p.SyncSend("topic1", "message1", func(err error){})

// 同步发送消息，指定key，相同key在同一分区
err := p.SendMessage(kafka.ProducerMessage{Topic: "topic1", Key: "key1", Value: "value1"}, func(err error){})
```
##### 关闭
```go
p.Close()
```
#### 消费者
##### 获取实例
```go
// 配置文件方式
c := kafka.GetConsumer("c")

// 手动实例化
opts := kafka.NewDefaultOptions()
opts.Addr = []string{"kafka:9092"}
opts.Consumer.Group = "ngo"
c, err := kafka.NewConsumer(opts)
```
##### 消费消息
```go
// 注册监听
c.AddListener("topic1", &listener{})
// 启动
c.Start()


type listener struct {
	kafka.Listener
}

func (l *listener) Listen(message kafka.ConsumerMessage, ack *kafka.Acknowledgment) {
    // 处理消息
    ...
    
    // 手动提交，当EnableAutoCommit=false时，下面代码才生效
	ack.Acknowledge()
}
```
##### 停止后台消费任务
```go
c.Stop()
```
### 注意事项
- 当kafka服务故障时，kafka生产端会自动熔断，但是系统默认配置比较宽松，*请务必配置好相关超时和重试参数*，否则生产时请求会非常慢，容易拖垮服务。

### 使用示例
- [examples/kafka](../examples/kafka)