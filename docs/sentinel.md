# [Ngo](https://github.com/NetEase-Media/ngo)

---
## sentinel
### 模块用途
sentinel模块源自阿里开源的sentinel，相关文档如下：
- https://github.com/alibaba/sentinel-golang/wiki
- https://sentinelguard.io/zh-cn/docs/golang/basic-api-usage.html

ngo做了简单的封装，集成了哨兵监控

### 使用说明
#### 配置
参数配置详见[sentinel options](config.md#sentinel-配置)

#### 调用

v0.1.57版本之前的配置将被弃用
```
circuitbreaker:
  rules:
    - resource: count
      strategy: 2
      retryTimeoutMs: 30
      minRequestAmount: 1
      statIntervalMs: 5000
      maxAllowedRtMs: 10
      threshold: 1.0
    - resource: ratio
      strategy: 1
      retryTimeoutMs: 30
      minRequestAmount: 10
      statIntervalMs: 5000
      maxAllowedRtMs: 10
      threshold: 1.0
```
新配置如下:
```
sentinel:
  circuitbreakerRules:
    - resource: count
      strategy: 2
      retryTimeoutMs: 30
      minRequestAmount: 1
      statIntervalMs: 5000
      maxAllowedRtMs: 10
      threshold: 1.0
    - resource: ratio
      strategy: 1
      retryTimeoutMs: 30
      minRequestAmount: 10
      statIntervalMs: 5000
      maxAllowedRtMs: 10
      threshold: 1.0
  flowRules:
  hotspotRules:
  isolationRules:
  systemRules:
```
规则参数描述，在官方文档中已经写的很清楚了，这里不在配置文件上加了。

```
import github.com/NetEase-Media/ngo/adapter/sentinel

sentinel.Entry("abc")
sentinel.TraceError(entry, err)

```

### 注意事项
- 业务调用请务必使用ngo包里的Entry和TraceError方法，否则哨兵无法监控