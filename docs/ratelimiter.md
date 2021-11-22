# [Ngo](https://github.com/NetEase-Media/ngo)

---
## 限流
### 模块用途
对http接口限流
### 使用说明
基于sentinel实现，流控相关参数见官网 [flow control](https://sentinelguard.io/zh-cn/docs/golang/flow-control.html)
该中间件可作用于全局、组、或者路由上

```go
RateLimiter(
	WithResource("abc"),    // 必填，关联sentinel资源名
    WithErrorHttpCode(429), // 选填，自定义错误响应码
    WithDefaultMsg("xxx"),  // 选填，自定义错误响应内容
    WithErrorHandler(errorHandler)  // 选填，自定义错误响应处理器，当此项存在，则忽略自定义响应码和内容
)
```


## 使用示例
- [examples/ratelimiter](../examples/ratelimiter)