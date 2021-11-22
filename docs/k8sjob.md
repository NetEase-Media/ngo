# [Ngo](https://github.com/NetEase-Media/ngo)

---
## k8s job

### 模块用途

在k8s pod中运行一次性的进程，从任务中心拉取分片参数，执行用户注册的函数，在完成后上报运行结果。

### 使用说明

job进程主要分为三步：

1. 任务中心为job按pod name和namespace分配分片id，job库后台发送请求获取。
2. job库将参数传给用户函数并执行。这里是用户唯一需要开发的地方，只需实现任务逻辑即可。 
3. 在用户函数执行完之后job库会上传执行时间和结果到任务中心。

### 使用示例

和其它服务一样启动需要有ngo的配置文件，在文件中加入job的参数：任务中心地址和namespace。

```yaml
job:
  centerUrl: http://127.0.0.1:9090
  namespace: maomao
```

在`main`函数中调用`job.Run`：

```go
job.Run(func(args *job.Args) (string, error) {
    return strconv.Itoa(args.SharedNum), nil
})
```
