# [Ngo](https://github.com/NetEase-Media/ngo)

---

## 启动

### 使用说明
go main
```go
package main

import (
	"context"
	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/adapter/protocol"
	"github.com/NetEase-Media/ngo/server"
	"github.com/gin-gonic/gin"
)

// go run . -c ./app.yaml
func main() {
	s := server.Init()
	s.PreStart = func() error {
		log.Info("do pre-start...")
		return nil
	}

	s.PreStop = func(ctx context.Context) error {
		log.Info("do pre-stop...")
		return nil
	}

	s.AddRoute(server.GET, "/hello", func(ctx *gin.Context) {
		ctx.JSON(protocol.JsonBody("hello"))
	})
	s.Start()
}
```
配置文件app.yaml
```yaml
service:
  appName: ngo-demo
  clusterName: ngo-demo-local

```
以上是最简单的启动例子，更多配置参考[yaml配置说明](config.md)，server中主要函数说明：
- `server.Init()`：初始化函数，在这里会初始化yaml中的配置
- `server.PreStart()`：启动前执行的函数，使用方可以在这里执行一些自定义的初始化
- `server.PreStop(ctx)`：关闭前执行的函数，使用方可以在这里执行一些自定义的资源释放操作
- `server.Stop(ctx)`：关闭函数，使用方无需主动调用，框架会监听信号来调用该函数，该函数用来关闭server以及中间件client等

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