# [Ngo](https://github.com/NetEase-Media/ngo)

---
## web中间件
### 模块用途
这里中间件为gin 中间件，提供web层的拦截、扩展，类似java中的filter，目前系统内置中间件如下：
* [accesslog](accesslog.md)，内置，yaml可配置参数
* [限流](ratelimiter.md)，使用方可直接使用
* [超时](timeout.md)，使用方可直接使用
* [admin auth](jwt-auth.md)，内置，yaml可配置参数
* [分号转换](semicolon.md)，内置，未对外暴露
### 使用说明
调用方式支持全局和路由方式，和Gin调用方式类似。
- global  
```go
func main() {
    s := server.Init()
    s.Use(timeout.Timeout(timeout.WithTimeout(50*time.Millisecond)))
    s.AddRoute(server.GET, "/hello", func(c *gin.Context) {
        c.String(http.StatusOK, "success")
    })
    s.Start()
}
```
- route
```go
func main() {
    s := server.Init()
	s.AddRoute(server.GET, "/hello", timeout.Timeout(timeout.WithTimeout(50*time.Millisecond)), func(c *gin.Context) {
        c.String(http.StatusOK, "success")
    })
    s.Start()
}
```