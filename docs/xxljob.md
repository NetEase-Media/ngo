# [Ngo](https://github.com/NetEase-Media/ngo)

---
## xxljob
### 模块用途
* 集成xxljob

### 使用说明
* 参考[配置说明](config.md)
* 配置好之后只用调用注册任务接口，之后在xxljob后台配置即可

### 注意事项
暂无

### 使用示例
```go
xxljob.RegTask("helloworld", func(cxt context.Context, param *xxl.RunReq, logger *XxlJobLogger) string {
  logger.Infof("123")
  logger.Errorf("Errorf test")
  return "ok"
})
```
