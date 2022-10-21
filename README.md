# [Ngo](https://github.com/NetEase-Media/ngo)

## 简介
Ngo是由网易传媒基础技术团队开源的一个类似Java Spring Boot的框架，全部使用Go语言开发，主要目标是：
- 提供比原有Java框架更高的性能和更低的资源占用率
- 尽量为业务开发者提供所需的全部工具库
- 自动加载配置和初始化程序环境，开发者能直接使用各种库
- 提供预置的健康检查、上线、下线、状态查看等接口

## 使用
使用 `go get -u github.com/NetEase-Media/ngo` 命令下载安装

## 快速开始
`main.go`
```go
package main

import (
	_ "github.com/NetEase-Media/ngo/pkg/include"
	"github.com/NetEase-Media/ngo/pkg/ngo"
	"github.com/NetEase-Media/ngo/pkg/server/http"
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
配置文件 `app.yaml`
```yaml
service:
  appName: ngo-demo
  clusterName: ngo-demo-local
httpServer:
  port: 8080

```


### 更多示例

- [examples](/examples)
- [ngo-demo](https://github.com/NetEase-Media/ngo-demo)

# 使用文档
* [web server](docs/server.md)
    * [启动](docs/start.md)
        * [yaml配置说明](docs/config.md)
        * [多环境配置&&配置导入](docs/yamlimport.md)
        * [pprof](docs/pprof.md)
    * [优雅停服](docs/gracefulshutdown.md)
    * [web中间件](docs/middleware.md)
        * [accesslog](docs/accesslog.md)
        * [限流](docs/ratelimiter.md)
        * [超时](docs/timeout.md)
        * [admin auth](docs/jwt-auth.md)
        * [分号转换](docs/semicolon.md)
* [日志](docs/log.md)
* [协议](docs/protocol.md)
* [sentinel](docs/sentinel.md)
* [中间件client](docs/client.md)
    * [db](docs/db.md)
    * [httplib](docs/httplib.md)
    * [kafka](docs/kafka.md)
    * [memcache](docs/memcache.md)
    * [多级缓存](docs/multicache.md)
    * [redis](docs/redis.md)
    * [分布式锁](docs/dlock.md)
    * [zookeeper](docs/zookeeper.md)
* [定时任务]()
    * [cron定时任务](docs/cron.md)
    * [k8s job](docs/k8sjob.md)
    * [xxljob](docs/xxljob.md)
* [工具](docs/util.md)


## 问题反馈
 对应bug上报、问题咨询和讨论，可以提交issue
