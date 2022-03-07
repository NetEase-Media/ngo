# [Ngo](https://github.com/NetEase-Media/ngo)

---

## 配置选项

### 配置示例
[configs/config_sample.yaml](../configs/config_sample.yaml)

### 配置说明
#### service 配置 (server.ServiceOptions)
| 字段名      | 类型   | 含义   | 必填 | 默认值                                   | 备注 |
| ---         | ---    | ---    | ---  | ---                                      | --   |
| appName     | string | 应用名 | 否   | 默认从环境变量`com_cmdb_appname`中取     |      |
| clusterName | string | 集群名 | 否   | 默认从环境变量`com_cmdb_clustername`中取 |      |

#### httpServer 配置 (server.Options)

**[gin文档](https://github.com/gin-gonic/gin)**

| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| port   | int    | 端口号   | 否   | 8080    |                                       |
| mode   | string | gin 模式 | 否   | release | 可选有 `["debug", "release", "test"]` |
| shutdownTimeout   | time.Duration | 停服超时时间 | 否   | 10s | |
| middlewares   | middleware struct 见下表 |中间件 | 否   |  | |

##### middleware 配置 (server.MiddlewaresOptions)
| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| accessLog   | AccessLogMwOptions struct  见下表  |   | 否   |     ||
| urlMetrics   | UrlMetricsMwOptions struct  见下表 |  |  否  |  |  |
| jwtAuth   | JwtAuthMwOptions struct  见下表  | | 否 |  | |

###### accessLog 配置 (server.AccessLogMwOptions)
| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| enabled   | bool  | 是否开启  | 否   | true    ||
| pattern   | string | 匹配格式 |  否  | %h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i" |  |
| path            | string            | 日志目录             | 是   | 空串   | 支持相对路径                                                                      |
| noFile          | bool              | 是否只显示到标准输出 | 否   | true   |                                                                                   |
| filePathPattern | string            | 文件名模式           | 否   | 空串   | 如果设置此值，则可自定义文件名模式                                                |
| maxAge          | time.Duration     | 日志保留时长         | 否   | 7天    | time.Duration可用单位有(ns, us, ms, s, m, h)，下同,例如`300ms, 2h45m`             |
| rotationTime    | time.Duration     | 日志滚动切割时长     | 否   | 24h    |                                                                                   |
| rotationSize    | int64             | 日志滚动切割体积     | 否   | 1024   | 单位为M                                                                           |

##### urlMetrics 配置 (server.UrlMetricsMwOptions)
| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| enabled   | bool  | 是否开启  | 否   | true    ||
| originalPath   | bool | 是否上报原始路径 |  否  | false | 对于restful的url上报原始路径需要在哨兵配置文件中配置聚合规则 |

##### jwtAuth 配置 (server.JwtAuthMwOptions)
| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| enabled   | bool  | 是否开启  | 否   | false    ||
| authHeader   | string  | 认证头  | 否   | Authorization    ||
| tokenType   | string  | token类型  | 否   | Bearer    ||
| accessTokenExpiresIn   | int  | 请求token失效时间  | 否   | 3600    ||
| refreshTokenExpiresIn   | int  | 刷新token失效时间  | 否   | 7200    ||
| oidc   | Oidc struct  | oidc 参数 | 否   |     ||
| routePathPrefix   | string  | 是否开启  | 否   | 空串    ||
| ignorePaths   | []string  | 忽略认证路径  | 否   | []    ||
###### oidc 配置 (oidc.Options)
| 字段名 | 类型   | 含义     | 必填 | 默认值  | 备注                                  |
| ---    | ---    | ---      | ---  | ---     | --                                    |
| clientId   | string  | 客户端id  | 是   |     ||
| clientSecret   | string  | 客户端密钥  | 是   |     ||
| encryption   | string  | 加密类型  | 否   | HS256    ||

#### log 配置 ([]log.Options)

**[logrus文档](https://github.com/sirupsen/logrus)**

**[file-rotatelogs文档](https://github.com/lestrrat-go/file-rotatelogs)**

| 字段名          | 类型              | 含义                 | 必填 | 默认值 | 备注                                                                              |
| ---             | ---               | ---                  | ---  | ---    | --                                                                                |
| name            | string            | logger名称           | 是   | 空串   | 如果不写且在数组中为第一个，那么其默认值为 default                                |
| level           | string            | 日志级别             | 否   | info   | 可选有 `["trace", "debug", "info","warn", "error", "fatal", "panic"]`             |
| packageLevel    | map[string]string | 包日志级别           | 否   | 空map  | key为包名称，value为日志级别，注意此处日志级别如果超过level定义的日志级别，则无效 |
| path            | string            | 日志目录             | 是   | 空串   | 支持相对路径                                                                      |
| errorPath       | string            | 错误日志路径         | 是   | 空串   | 支持相对路径                                                                      |
| fileName        | string            | 日志文件名           | 是   | 空串   | 如果不写且在数组中为第一个，那么其默认值为 service.appName                        |
| writableStack   | bool              | 是否打印错误堆栈     | 否   | false  |                                                                                   |
| format          | string            | 日志格式             | 否   | txt    | 可选有 `["json", "txt", "blank" //表示空白]`                                      |
| noFile          | bool              | 是否只显示到标准输出 | 否   | true   |                                                                                   |
| filePathPattern | string            | 文件名模式           | 否   | 空串   | 如果设置此值，则可自定义文件名模式                                                |
| maxAge          | time.Duration     | 日志保留时长         | 否   | 7天    | time.Duration可用单位有(ns, us, ms, s, m, h)，下同,例如`300ms, 2h45m`             |
| rotationTime    | time.Duration     | 日志滚动切割时长     | 否   | 24h    |                                                                                   |
| rotationSize    | int64             | 日志滚动切割体积     | 否   | 1024   | 单位为M                                                                           |

##### level packageLevel 说明
* level 设定当前logger的日志级别，日志级别以及其对应值为 `{"panic":0, "fatal":1, "error":2, "warn":3, "info":4, "debug":5, "trace":6}`
* 根据以上值定义，如果level设定为info，那么5,6对应的debug和trace级别都不会被打印
* packageLevel 设定某个包日志级别，这个日志级别是局限在level的设定下的
* 假如设置 `{"gorm.io/gorm": "error"}`，那么 gorm.io/gorm 这个包下的 error以上的日志都不会被打印
* 假如设置 `{"gorm.io/gorm": "debug"}`， 因为其受限于level不能实现，所以debug级别的日志不能被打印，其相当于日志级别设置为info

#### redis 配置 ([]redis.Options)

**[goredis文档](https://github.com/go-redis/redis)**

| 字段名             | 类型          | 含义                               | 必填                                     | 默认值              | 备注                                                         |
| ---                | ---           | ---                                | ---                                      | ---                 | --                                                           |
| name               | string        | redis配置名称                      | 是                                       | 空串                | 需唯一                                                       |
| connType           | string        | 连接类型                           | 是                                       | 空串                | 可选有 `["client", "cluster", "sentinel", "sharded_sentinel" |
| addr               | []string      | 地址列表                           | 是                                       | 无                  | 格式为 `host:port` 如果是单实例只会取第一个                  |
| username           | string        | 用户名                             | 是                                       | 空串                |                                                              |
| password           | string        | 密码                               | 是                                       | 空串                |                                                              |
| masterNames        | []string      | master名称列表                     | 只有sentinel sharded_sentinel 类型时必填 | nil                 |                                                              |
| autoGenShardName   | bool          | 是否自动生成分片名称               | 否                                       | false               | 只有sharded_sentinel 类型时有用      |
| db                 | int           | 数据库                             | 否                                       | 0                   |                                                              |
| maxRetries         | int           | 最大重试次数                       | 否                                       | 3                   | 传-1表示禁用重试                                             |
| minRetryBackoff    | time.Duration | 最小重试backoff                    | 否                                       | 8ms                 | 传-1表示禁用backoff                                          |
| maxRetryBackoff    | time.Duration | 最大重试backoff                    | 否                                       | 512ms               | 传-1表示禁用backoff                                          |
| dialTimeout        | time.Duration | 连接超时时间                       | 否                                       | 5s                  |                                                              |
| readTimeout        | time.Duration | 读超时                             | 否                                       | 3s                  | 传-1表示不超时                                               |
| writeTimeout       | time.Duration | 写超时                             | 否                                       | 3s                  | 传-1表示不超时                                               |
| poolSize           | int           | 最大连接数                         | 否                                       | 10 * gomaxproces(0) |                                                              |
| minIdleConns       | int           | 最小空闲连接数                     | 否                                       | 0                   |                                                              |
| maxConnAge         | time.Duration | 最大连接时间                       | 否                                       | 默认不关闭连接      |                                                              |
| poolTimeout        | time.Duration | 所有连接都在工作是需要等待超时时间 | 否                                       | readTimeout+1s      |                                                              |
| idleTimeout        | time.Duration | 关闭空闲连接需要等待的时间         | 否                                       | 5m                  |                                                              |
| idleCheckFrequency | time.Duration | 空闲连接检查时间间隔               | 否                                       | 1m                  |                                                              |

##### autoGenShardName 说明
* 如果为false，默认使用MasterName， 只当sharded_sentinel 类型使用。
* 该字段用来兼容旧项目，非特殊情况请勿设置成true，否则在MasterNames顺序变化时会造成分配rehash

#### memcache 配置 ([]memcache.Options)

**[gomemcache文档](https://github.com/bradfitz/gomemcache)**

| 字段名       | 类型          | 含义             | 必填 | 默认值      | 备注               |
| ---          | ---           | ---              | ---  | ---         | --                 |
| name         | string        | memcache配置名称 | 是   | 空串        | 需唯一             |
| addr         | []string      | 地址列表         | 是   | nil         | 格式为 `host:port` |
| timeout      | time.Duration | 连接超时时间     | 否   | 100ms       |                    |
| maxIdleConns | int           | 最大空闲连接数   | 否   | 每个地址2个 |                    |

####  db 配置 ([]db.Options)

**[gorm文档](https://gorm.io/docs/)**

| 字段名          | 类型          | 含义             | 必填 | 默认值 | 备注                                                                         |
| ---             | ---           | ---              | ---  | ---    | --                                                                           |
| name            | string        | db配置名称       | 是   | 空串   | 需唯一                                                                       |
| type            | string        | db类型       | 否   | mysql   | 可选有 ["mysql", "ddb"]                                                                       |
| url             | string        | db连接地址       | 是   | 空串   | [连接串格式](https://colobu.com/2019/01/10/drivers-connection-string-in-Go/) |
| maxIdleConns    | int           | 最大空闲连接数   | 否   | 2      |                                                                              |
| maxOpenConns    | int           | 最大连接数       | 否   | 0      | 0为无限制                                                                    |
| connMaxLifetime | time.Duration | 连接最长存活时间 | 否   | 0      | 小于等于0 连接将一直存在                                                     |
| connMaxIdleTime | time.Duration | 连接最长空闲时间 | 否   | 0      | 小于等于0 连接将一直存在                                                     |

*注意：如果是ddb，确定是否能使用服务端预处理，如果不能请在url参数中设置InterpolateParams=true*

####  kafka 配置 ([]kafka.Options)

**[sarama文档](https://github.com/Shopify/sarama)**

**[sarama/config.go](https://github.com/Shopify/sarama/blob/master/config.go)**

| 字段名       | 类型                   | 含义          | 必填 | 默认值 | 备注               |
| ---          | ---                    | ---           | ---  | ---    | --                 |
| name         | string                 | kafka配置名称 | 是   | 空串   | 需唯一             |
| addr         | []string               | 地址列表      | 是   | nil    | 格式为 `host:port` |
| version      | string                 | 版本          | 否   | 2.1.0  |                    |
| maxOpenRequests  | int          | 最大打开连接数      | 否   | 5    |                    |
| dialTimeout  | time.Duration          | 连接超时      | 否   | 30s    |                    |
| readTimeout  | time.Duration          | 读超时        | 否   | 30s    |                    |
| writeTimeout | time.Duration          | 写超时        | 否   | 30s    |                    |
| metadata     | metadata struct 见下表 |               |      |        |                    |
| consumer     | consumer struct 见下表 |               |      |        |                    |
| producer     | producer struct 见下表 |               |      |        |                    |

##### kafka metadata 配置

**[sarama/consumer.go](https://github.com/Shopify/sarama/blob/master/consumer.go)**

| 字段名             | 类型          | 含义                 | 必填 | 默认值     | 备注                                                                       |
| ---                | ---           | ---                  | ---  | ---        | --                                                                         |
| retries            | int           | 重试次数             | 否   | 3          |                                                                            |
| timeout            | time.Duration | 会话超时             | 否   | 60s          |                                                                            |

##### kafka consumer 配置

**[sarama/consumer.go](https://github.com/Shopify/sarama/blob/master/consumer.go)**

| 字段名             | 类型          | 含义                 | 必填 | 默认值     | 备注                                                                       |
| ---                | ---           | ---                  | ---  | ---        | --                                                                         |
| group              | string        | 消费组               | 是   |         |                                                                            |
| enableAutoCommit   | bool          | 是否自动提交offset   | 否   | true       |                                                                            |
| autoCommitInterval | time.Duration | 自动提交时间间隔     | 否   | 1s         |                                                                            |
| initialOffset      | int64         | 消费的初始offset     | 否   | 最新offset | 这个配置在在之前没有提交过offset的情况下有效                               |
| sessionTimeout     | time.Duration | 会话超时             | 否   | 10s        | 类似java中`group.min.session.timeout.ms` 和 `group.max.session.timeout.ms` |
| minFetchBytes      | int32         | 每次请求最小字节数   | 否   | 1          | 等价于java中的`fetch.min.bytes`                                            |
| defaultFetchBytes  | int32         | 每次请求默认字节数   | 否   | 1024 * 1024| 类似java中的`fetch.message.max.bytes`                                            |
| maxFetchBytes      | int32         | 每次请求最大字节数   | 否   | 0          | 类似java中的`fetch.message.max.bytes`,0表示无限制                          |
| maxFetchWait       | time.Duration | 每次请求最长等待时间 | 否   | 250ms      | 类似java中的`fetch.wait.max.ms`                                            |
| retries            | int           | 重试次数             | 否   | 3          |                                                                            |

##### kafka producer 配置

**[sarama/sync_producer.go](https://github.com/Shopify/sarama/blob/master/sync_producer.go)**

**[sarama/async_producer.go](https://github.com/Shopify/sarama/blob/master/async_producer.go)**

| 字段名           | 类型          | 含义                | 必填 | 默认值      | 备注                                                                                                                         |
| ---              | ---           | ---                 | ---  | ---         | --                                                                                                                           |
| maxMessageBytes  | int           | 消息最大字节数      | 否   | 1000000     |                                                                                                                              |
| acks             | int16         | ack机值             | 否   | 1           | 可选有`{0:"NoResponse,不等ACK", 1:"WaitForLocal,本地ack":1,-1:"WaitForLocal,副本ack"}` 等价于java中的`request.required.acks` |
| timeout          | time.Duration | 等待ack的超时时间   | 否   | 10秒        | 等价于java中的`request.timeout.ms`                                                                                           |
| retries          | int           | 重试次数            | 否   | 3           | 类似java中的`message.send.max.retries`                                                                                       |
| maxFlushBytes    | int           | 触发flush的字节大小 | 否   | 100*1048576 |                                                                                                                              |
| maxFlushMessages | int           | 触发flush的消息大小 | 否   | 0           |                                                                                                                              |
| flushFrequency   | time.Duration | 触发flush的时间间隔 | 否   | 1s         | 类似java中的`queue.buffering.max.ms`                                                                                         |
| idempotent   | bool | 发送是否幂等 | 否   | false         | 如果为true，必须要求acks = -1                                                                                      |

####  sentinel 配置
| 字段名              | 类型                   | 含义             | 必填 | 默认值 | 备注                                                                                       |
| ---                 | ---                    | ---              | ---  | ---    | --                                                                                         |
| circuitBreakerRules | []*circuitbreaker.Rule | 熔断规则         | 否   | 无     | [熔断文档](https://sentinelguard.io/zh-cn/docs/golang/circuit-breaking.html)               |
| flowRules           | []*flow.Rule           | 流控规则         | 否   | 无     | [流控文档](https://sentinelguard.io/zh-cn/docs/golang/flow-control.html)                   |
| hotspotRules        | []*hotspot.Rule        | 热点参数流控规则 | 否   | 无     | [热点文档](https://sentinelguard.io/zh-cn/docs/golang/hotspot-param-flow-control.html)     |
| isolationRules      | []*isolation.Rule      | 流量隔离规则     | 否   | 无     | [隔离文档](https://sentinelguard.io/zh-cn/docs/golang/concurrency-limiting-isolation.html) |
| systemRules         | []*system.Rule         | 自适应流控规则   | 否   | 无     | [自适应文档](https://sentinelguard.io/zh-cn/docs/golang/system-adaptive-protection.html)   |

##### circuitBreakerRules 配置 (circuitbreaker.Rule)

**[熔断降级](https://sentinelguard.io/zh-cn/docs/golang/circuit-breaking.html)**

| 字段名           | 类型    | 含义               | 必填 | 默认值 | 备注                                                                             |
| ---              | ---     | ---                | ---  | ---    | --                                                                               |
| resource         | string  | 熔断器埋点资源名称 | 是   | 空串   |                                                                                  |
| strategy         | uint32  | 熔断策略           | 是   | 0      | 可选有`{0: "SlowRequestRatio", 1: "ErrorRatio", 2: "ErrorCount"}`,详情见相关文档 |
| retryTimeoutMs   | uint32  | 熔断触发后持续时间 | 是   | 0      | 单位为毫秒                                                                       |
| statIntervalMs   | uint32  | 统计的时间窗口长度 | 是   | 0      | 单位为毫秒                                                                       |
| minRequestAmount | uint64  | 静默数量           | 否   | 0      | 触发熔断的最小请求数目                                                           |
| maxAllowedRtMs   | uint64  | 请求慢调用的临界值 | 否   | 0      | 单位为毫秒                                                                       |
| threshold        | float64 | 慢调用比例的阈值   | 否   | 0      | 仅SlowRequestRatio模式有效，小数表示，比如0.1表示10%                             |

##### flowRules 配置 (flow.Rule)

**[流量控制](https://sentinelguard.io/zh-cn/docs/golang/flow-control.html)**

| 字段名                 | 类型    | 含义                                         | 必填 | 默认值 | 备注                                                                |
| ---                    | ---     | ---                                          | ---  | ---    | --                                                                  |
| resource               | string  | 流控资源名称                                 | 是   | 空串   |                                                                     |
| tokenCalculateStrategy | int32   | 流控token计算策略                            | 是   | 0      | 可选有 `{0:"Direct", 1:"WarmUp", 2:"MemoryAdaptive"}`               |
| controlBehavior        | int32   | 控制策略                                     | 是   | 0      | 可选有 `{0:"Reject,超过阈值直接拒绝", 1:"Throttling,匀速排队"}      |
| threshold              | float64 | 流控阈值                                     | 是   | 0      | 如果字段`StatIntervalInMs`是1000(也就是1秒)，那么Threshold就表示QPS |
| statIntervalInMs       | uint32  | 规则对应的流量控制器的独立统计结构的统计周期 | 是   | 0      | 如果StatIntervalInMs是1000，也就是统计QPS                           |
| relationStrategy       | int32   | 调用关系限流策略                             | 否   | 0      | 可选有`{0:"CurrentResource", 1:"AssociatedResource"}                |
| refResource            | string  | 关联的resource                               | 否   | 空串   |                                                                     |
| warmUpPeriodSec        | uint32  | 预热的时间长度                               | 否   | 0      | 仅仅对`WarmUp`的`TokenCalculateStrategy`生效                        |
| warmUpColdFactor       | uint32  | 预热的因子                                   | 否   | 3      | 仅仅对`WarmUp`的`TokenCalculateStrategy`生效                        |
| maxQueueingTimeMs      | uint32  | 匀速排队的最大等待时间                       | 否   | 0      | 仅对 `Throttling ControlBehavior`生效                               |

##### hotspotRule 配置 (hotspot.Rule)

**[热点参数流控](https://sentinelguard.io/zh-cn/docs/golang/hotspot-param-flow-control.html)**

| 字段名            | 类型                  | 含义                              | 必填 | 默认值 | 备注                                                                |
| ---               | ---                   | ---                               | ---  | ---    | --                                                                  |
| resource          | string                | 资源名                            | 是   | 空串   |                                                                     |
| metricType        | int32                 | 流控指标类型                      | 是   | 0      | 可选有`{0:"Concurrency,并发", 1:"QPS"}`                             |
| controlBehavior   | int32                 | 流控的效果                        | 是   | 0      | 可选有`{0:"Reject, 超过阈值直接拒绝", 1:"Throttling,匀速排队"}`     |
| paramIndex        | int                   | 热点参数的索引                    | 是   | 0      | 对应 WithArgs(args ...interface{}) 中的参数索引位置，从 0 开始      |
| threshold         | int64                 | 限流阈值                          | 是   | 0      | 针对每个热点参数                                                    |
| maxQueueingTimeMs | int64                 | 最大排队等待时长                  | 否   | 0      | 仅在匀速排队模式 + QPS 下生效                                       |
| burstCount        | int64                 | 静默值                            | 否   | 0      | 仅在快速失败模式 + QPS 下生效                                       |
| durationInSec     | int64                 | 统计结构填充新的 token 的时间间隔 | 否   | 0      | 仅在请求数(QPS)流控模式下生效                                       |
| paramsMaxCapacity | int64                 | 统计结构的容量最大值              | 否   | 20000  | Top N                                                               |
| specificItems     | map[interface{}]int64 | 特定参数的特殊阈值配置            | 否   | 空map  | 可以针对指定的参数值单独设置限流阈值，不受前面 Threshold 阈值的限制 |

##### isolationRules 配置 (isolation.Rule)

**[并发隔离控制](https://sentinelguard.io/zh-cn/docs/golang/concurrency-limiting-isolation.html)**

| 字段名     | 类型   | 含义     | 必填 | 默认值 | 备注                                                         |
| ---        | ---    | ---      | ---  | ---    | --                                                           |
| resource   | string | 资源名   | 是   | 空串   |                                                              |
| metricType | int32  | 指标类型 | 是   | 0      | 可选有`{0:"Concurrency,并发"}`                               |
| threshold  | uint32 | 隔离阈值 | 是   | 0      | 如果资源当前的并发数高于阈值 (Threshold)，那么资源将不可访问 |

##### systemRules 配置 (system.Rule)

**[系统自适应流控](https://sentinelguard.io/zh-cn/docs/golang/system-adaptive-protection.html)**

| 字段名       | 类型    | 含义         | 必填 | 默认值 | 备注                                                                                                                   |
| ---          | ---     | ---          | ---  | ---    | --                                                                                                                     |
| metricType   | int32   | 指标类型     | 是   | 0      | 可选有`{0:"Load,load1负载", 1:"AvgRT,平均响应时间", 2:"Concurrency,并发", 3:"InboundQPS,QPS", 4:"CpuUsage,CPU使用率"}` |
| triggerCount | float64 | 最低临界值   | 是   | 0      |                                                                                                                        |
| strategy     | int32   | 自适应性策略 | 是   | 0      | 可选有`{-1:"NoAdaptive,无策略", 0:"BBR,类似TCP BBR"}`                                                                  |

####  httpClient 配置 (httplib.Options)

**[fasthttp文档](https://pkg.go.dev/github.com/valyala/fasthttp)**

| 字段名                    | 类型          | 含义                             | 必填 | 默认值 | 备注                                |
| ---                       | ---           | ---                              | ---  | ---    | --                                  |
| name                      | string        | 客户端名称 User-Agent header的值 | 否   | 空串   |                                     |
| noDefaultUserAgentHeader  | bool          | 是否设置User-Agent header        | 否   | fasle  |                                     |
| maxConnsPerHost           | int           | 每个host的最大连接数             | 否   | 512    |                                     |
| maxIdleConnDuration       | time.Duration | 空闲的keep-alive连接最大关闭时间 | 否   | 10s    |                                     |
| maxConnDuration           | time.Duration | keep-alive连接最大关闭时间       | 否   | 无限制 |                                     |
| maxIdemponentCallAttempts | int           | 最大重试次数                     | 否   | 5      |                                     |
| readBufferSize            | int           | 每个连接的读缓存大小             | 否   | 4096   |                                     |
| writeBufferSize           | int           | 每个连接的写缓存大小             | 否   | 4096   |                                     |
| readTimeout               | time.Duration | 读超时                           | 是   | 无限制 | 建议设置                            |
| writeTimeout              | time.Duration | 写超时                           | 是   | 无限制 | 建议设置                            |
| maxResponseBodySize       | int           | 最大响应体大小                   | 否   | 无限制 |                                     |
| maxConnWaitTimeout        | time.Duration | 等待空闲连接的最大时间           | 否   | 不等待 | 默认不等待，直接返回 ErrNoFreeConns |

####  tracer 配置 (tracer.Options)

| 字段名  | 类型   | 含义               | 必填 | 默认值          | 备注 |
| ---     | ---    | ---                | ---  | ---             | --   |
| updHost | string | tracer上报udp host | 否   | 160.254.169.253 |      |
| updPort | string | tracer上报udp port | 否   | 6831            |      |

####  dlock 配置 (dlock.Options)
| 字段名 | 类型     | 含义            | 必填 | 默认值 | 备注 |
| ---    | ---      | ---             | ---  | ---    | --   |
| pools  | []string | redis配置名列表 | 是   | 空串   |      |

####  pprof 配置 (server.PprofOptions)

**[pprof文档](https://golang.org/pkg/net/http/pprof/)**

| 字段名 | 类型 | 含义          | 必填 | 默认值 | 备注 |
| ---    | ---  | ---           | ---  | ---    | --   |
| switch | bool | 是否开启pprof | 否   | false  |      |
| port   | int  | 端口号        | 是   | nil    |      |


####  miner 缓存配置 ([]multicache.Options)

**[gcache文档](https://github.com/bluele/gcache)**
**[goredis文档](https://github.com/go-redis/redis)**

| 字段名       | 类型   | 含义                  | 必填 | 默认值    | 备注                                                      |
| ---          | ---    | ---                   | ---  | ---       | --                                                        |
| type         | string | 缓存类型              | 是   | 空串      | 可选有`["local", "redis"]`                                |
| priority     | int    | 加载优先级            | 是   | nil       | 0为最高优先级，数字越大优先级越低, local默认0, redis默认1 |
| capacity     | int    | 最大容量，用于local   | 否   | 1,000,000 |                                                           |
| defaultRedis | string | 默认使用的redis配置名 | 否   | 空串      | type为redis情况下必须填写                                 |


####  xxljob 配置 (xxljob.Options)

**[xxljob文档](https://www.xuxueli.com/xxl-job)**
**[xxl-job-executor-go 文档](https://github.com/xxl-job/xxl-job-executor-go)**

| 字段名       | 类型   | 含义           | 必填 | 默认值           | 备注                    |
| ---          | ---    | ---            | ---  | ---              | --                      |
| enabled      | bool   | 是否开启       | 是   | 否               | 是否开启xxljob          |
| serverAddr   | string | xxljob后台地址 | 是   | 空串             |                         |
| accessToken  | string | 请求令牌       | 否   | 空串             | 请求令牌,暂时可不用填写 |
| executorIp   | string | 本地执行器ip   | 否   | 默认取本机ip地址 |                         |
| executorPort | int    | 本地执行器端口 | 否   | 19876            |                         |
| registryKey  | string | 执行器名称     | 否   | 默认值为集群名   |                         |
| logDir       | string | 日志目录       | 否   | "./log/xxljob"   |                         |

































