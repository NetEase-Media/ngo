# [Ngo](https://github.com/NetEase-Media/ngo)

---

## zookeeper client

### 模块用途

提供Zookeeper 客户端。客户端是在[go-zookeeper](https://github.com/go-zookeeper/zk)基础上实现，由于原始接口比较复杂，业务需求一般用不上，Ngo中对其进行了较多的封装。在配置文件中增加zookeeper字段，Ngo会自动指定名称的zookeeper客户端。

### 使用说明

#### 配置

参数配置详见[zookeeper client](config.md#zookeeper-配置zookeeperoptions)

#### 各模式应用

##### 获取客户端

```go
// 配置文件方式
client := zookeeper.GetZkClient("zookeeperClientName")

// 手动实例化
var opts []zookeeper.Options
o := Options{
	Name: "zookeeperName",
	Addr: []string{"host:port"},
	SessionTimeout: 5*time.Second,
}
opts = append(opts,o)
zookeeperClients,err := zookeeper.NewClientFromZookeeper(opts)
client := zookeeperClients["zookeeperClientName"]
```

##### 执行命令

```go
// 创建节点(默认拥有全部节点权限)
	// flags有4种取值:
	// PERSISTENT                   // 持久化节点
	// EPHEMERAL                    // 临时节点， 客户端session超时这类节点就会被自动删除
	// PERSISTENT_SEQUENTIAL        // 顺序自动编号持久化节点，这种节点会根据当前已存在的节点数自动加 1
	// EPHEMERAL_SEQUENTIAL         // 临时自动编号节点
path,err := client.CreateNode(path string, flags int, data string)

// 创建节点(自定义节点权限)
path,err := client.CreateNodeWithAcls(path string, flags int32, acls []zk.ACL, data string)
	

// 创建节点（包含父节点）  如：想创建 /a/c 节点，但/a节点不存在，则把/a 和 /a/c 都建立
path,err := client.CreateNodeParent(path string, flags int, data string)

// 设置节点值
err := client.SetData(path string, data string)

// 获取节点值
data,err := client.GetData(path string)

// 获取子节点列表
pathSlices,err := client.GetChildren(path string)

// 删除节点
err := client.Delete(path string)

// 判断节点是否存在
isExist := client.Exist(path string)

// 获取zookeeper当前的连接状态
connState := client.GetConnState()

// 获取当前连接的sessionId 
sessionId := client.GetSessionId()

// 开启连接状态监听
ch,err := client.StartLi()

// 监听子节点                                注： 目前只能监听子节点的 增加 和  删除 ， 子节点值的改变 无法监听
client.WatchChildren(realPath, func(respChan <-chan *zookeeper.WatchChildrenResponse) {
	for r := range respChan {
		for k, _ := range r.ChildrenChangeInfo {
			fmt.Println(r.ChildrenChangeInfo[k].Path + "  " + r.ChildrenChangeInfo[k].OperateType.String())
		}
	}
})

// 监听当前节点
client.WatchNode(realPath, func(respChan <-chan *zookeeper.WatchNodeResponse) {
	for r := range respChan {
		fmt.Println("operate: " + r.NodeChangeInfo.OperateType.String() + " olddata: " + r.NodeChangeInfo.OldData + " newdata: " + r.NodeChangeInfo.NewData + " path: " + r.NodeChangeInfo.Path)
	}
})
```

命令范例详见 [zookeeper 命令范例](../client/zookeeper/zookeeper_test.go)

##### 关闭
```go
zookeeper.CloseAll()
```

### 使用示例

- [examples/zookeeper](../examples/zookeeper)
