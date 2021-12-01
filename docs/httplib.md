# [Ngo](https://github.com/NetEase-Media/ngo)

---

## httplib

### 模块用途

提供HTTP客户端及相关工具函数，来发送HTTP请求和解析回复。

### 使用说明

#### 客户端配置

ngo启动时会使用配置文件来初始化全局默认的HTTP客户端，如果需要多个不同配置的客户端可使用`New`来新建。

#### HTTP方法

当前支持GET、POST、PUT、DELETE方法。当调用几个方法接口如`httplib.Get()`时，会生成一个`DataFlow`对象，用来临时存储当前请求的各种中间变量。

#### 写入请求

可以使用接口写入以下字段：

- header
- ContentType（特殊header）
- query
- body（可以选择写入[]byte、对象json序列化、x-www-form-urlencoded）

#### 解析回复

使用bind系列函数可以设置回复返回后需要解析的格式，除`BindHeader`用来解析header，其它均为body接口，当前支持：

- 整型
- 浮点型
- 字符串
- byte数字
- 对象json序列化

#### 其它特殊功能

- 超时设置
- 服务降级
- 服务熔断

### 注意事项

注意如果使用了httplib最后一定要调用`Do`函数。

每次生成`DataFlow`对象时都会从`sync.Pool`中取出一个请求对象，直到调用`Do`后放回Pool。如果在代码中忘记调用`Do`会增加垃圾回收负担影响性能。

同样如果调用了熔断接口也会在`Do`中释放资源，如果未调用`Do`会造成熔断逻辑错误。

### 使用示例

Get请求

```go
httplib.Get("xxx").Do(ctx)
```

Post请求

```go
httplib.Post("xxx").SetJson(obj).Do(ctx)
```

绑定请求参数和回复

```go
httplib.Get("xxx").SetHeader(h).SetQuery(q).BindHeader(rh).BindInt(&i).Do(ctx)
```

服务降级

```go
httplib.Get("xxx").Degrade(func () error {
return nil
}).Do(ctx)
```

更多示例可见代码[example_test.go](../client/httplib/example_test.go)
