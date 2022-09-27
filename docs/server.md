# [Ngo](https://github.com/NetEase-Media/ngo)

---
## web server

### 模块说明
ngo web server 集成了httpserver（基于gin），pprof server，哨兵以及一些常用的组件，启动需要配合yaml配置文件。

* [启动](start.md)
    * [yaml配置说明](config.md)
    * [多环境yaml导入](yamlimport.md)
    * [pprof](pprof.md)
* [优雅停服](gracefulshutdown.md)
* [web中间件](middleware.md)
    * [accesslog](accesslog.md)
    * [限流](ratelimiter.md)
    * [超时](timeout.md)
    * [admin auth](jwt-auth.md)
    * [分号转换](semicolon.md)
