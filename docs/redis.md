# [Ngo](https://github.com/NetEase-Media/ngo)

---
## redis client
### 模块用途
提供Redis客户端，Redis客户端选择 [go-redis](https://github.com/go-redis/redis) 实现。同样只需在配置中提供Redis服务配置，即可在运行中直接使用GetClient获取指定名字的客户端。其支持client、cluster、sentinel、sharded sentinel四种模式的Redis连接，且都能自动上报哨兵监控数据。
### 使用说明
#### 配置
参数配置详见[redis options](config.md#redis-配置-redisoptions)
#### 各模式应用
##### 获取实例
```go
// 配置文件方式（各模式通用）
c = redis.GetClient("redis")

// 手动实例化（client 单机模式）
opt := redis.NewDefaultOptions()
opt.Name = "client1"
opt.ConnType = "client"
// 只取第一个
opt.Addr = []string{"client1:6379"}
opt.Username = "xxxx"
opt.Password = "xxxx"
c = redis.NewClient(opt)

// 手动实例化（sentinel 哨兵模式）
opt := redis.NewDefaultOptions()
opt.Name = "sentinel1"
opt.ConnType = redis.RedisTypeSentinel
opt.Addr = []string{"s1:9000", "s2:9000", "s3:9000"}
// 只取第一个
opt.MasterNames = []string{"xxxx"}
opt.Username = "xxxx"
opt.Password = "xxxx"
c = redis.NewSentinelClient(opt)

// 手动实例化（sharded_sentinel 哨兵分片模式）
opt := redis.NewDefaultOptions()
opt.Name = "sharded-sentinel1"
opt.ConnType = redis.RedisTypeShardedSentinel
opt.Addr = []string{"s1:9000", "s2:9000", "s3:9000"}
opt.MasterNames = []string{"xxxx", "xxxx"}
opt.Username = "xxxx"
opt.Password = "xxxx"
c = redis.NewShardedSentinelClient(opt)

// 手动实例化（cluster 集群模式）
opt := redis.NewDefaultOptions()
opt.Name = "cluster1"
opt.ConnType = redis.RedisTypeShardedSentinel
opt.Addr = []string{"c1:6379", "c2:6379", "c3:6379"}
opt.Username = "xxxx"
opt.Password = "xxxx"
c = redis.NewClusterClient(opt)

```
##### 执行命令
```go
// SET
val, err := c.Set(ctx, "key1", "value1", time.Second).Result()
// GET
val, err := c.Get(context.Background(), "key1").Result()
```
更多命令详见 [go-redis commands](https://github.com/go-redis/redis/blob/master/commands.go#L81) 
##### 关闭
```go
c.Close()
```

### 使用示例
- [examples/redis](../examples/redis) 