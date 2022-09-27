# [Ngo](https://github.com/NetEase-Media/ngo)

---

## 启动

### 使用说明
go main
```go
package main

import (
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	_ "github.com/NetEase-Media/ngo/pkg/server/http/admin"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	app := ngo.Init()
	s := http.Get()
	s.AddRoute(http.GET, "/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world!")
	})
	app.Start()
}
```
配置文件app.yaml
```yaml
service:
  appName: ngo-demo
  clusterName: ngo-demo-local
httpServer:
  port: 8080
```
以上是最简单的启动例子，更多配置参考[yaml配置说明](config.md)，server中主要函数说明：
- `ngo.Init()`：初始化函数，在这里会初始化yaml中的配置
- `app.PreStart()`：启动前执行的函数，使用方可以在这里执行一些自定义的初始化
- `app.AfterStop()`：关闭后执行的函数，使用方可以在这里执行一些自定义的资源释放操作
- `app.Stop()`：关闭函数，使用方无需主动调用，框架会监听信号来调用该函数，该函数用来关闭server以及中间件client等

内置路由：
- `/health/online`：流量灰度中容器上线时调用，允许服务开始接受请求
- `/health/offline`：流量灰度中容器下线时调用，停止接收请求
- `/health/check`：提供k8s liveness探针，展示当前进程存活状态
- `/health/status`：提供k8s readiness探针，表明当前服务状态，是否能提供服务

### 使用示例
- [examples/quickstart](../examples/quickstart)

### 更多
* [yaml配置说明](config.md)
* [多环境yaml导入](yamlimport.md)
* [pprof](pprof.md)