# [Ngo](https://github.com/NetEase-Media/ngo)

---
## 超时
### 模块用途
设置接口超时时间，超时后可自定义返回
### 使用说明
~~该中间件可作用于全局、组、或者路由上~~
该中间件只能作用于路由上

```go
timeout.Timeout(    
    timeout.WithTimeout(2*time.Second), // 必填，超时时间
    timeout.WithHandler(func(c *gin.Context) {           // 必填，业务处理逻辑
        c.JSON(http.StatusServiceUnavailable, gin.H{"code": -1})
    }),
    timeout.WithErrorHttpCode(http.StatusServiceUnavailable), // 选填，自定义错误响应码
    timeout.WithDefaultMsg(defaultMsg),                       // 选填，自定义错误响应内容
    timeout.WithErrorHandler(func(c *gin.Context) {           // 选填，自定义错误响应处理器，当此项存在，则忽略自定义响应码和内容
        c.JSON(http.StatusServiceUnavailable, gin.H{"code": -1})
    }),
	timeout.WithCallBack(func(r *http.Request) {              // 选填，超时回调，可用于打日志
		fmt.Println("timeout happen, url:", r.URL.String())
	})
)   
```

### 使用示例
- [examples/timeout](../examples/timeout)