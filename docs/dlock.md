# [Ngo](https://github.com/NetEase-Media/ngo)

---
## 分布式锁
### 模块用途
提供分布式锁能力，基于redis实现
### 使用说明
#### 配置说明

分布式锁依赖ngo redis，pools参数为redis配置中`name`的值，至少填一个名称，为满足高可用，可自行添加多个redis实例，
当len(pools)/2 + 1个client上锁成功，即获取到锁
``` 
dlock:
  pools:
    - client1
    - client2
```

#### 调用示例
```
succ, executed, err := dlock.NewMutex("test", func() {
		log.Info("start working...")
		// do something
		
		log.Info("end working...")
	}).DoContext(ctx)
```
#### 方法说明
DoContext 方法默认完成加锁、解锁以及快要超时的续租工作，我们只需要完成`func()`里的逻辑即可

注意：NewMutex默认尝试次数tries=32，如果仅执行一次，例如job，请修改 `dlock.NewMutex(...).WithTries(1).DoContext(ctx)`

超过重试次数会自动失败，可通过`WithRetryDelay`或者`WithRetryDelayFunc`修改重试策略
其他参数一般不用关心，如需设置，可通过WithXXX进行更改

#### 返回参数说明
- succ 整个调用是否成功
- executed 业务方法是否执行，出错可能在获取锁阶段，那么 executed=false
- err 错误描述

#### 错误组合描述
- false, false, nil 获取锁阶段，表示存在锁竞争，未获取锁
- false, false, not nil 获取锁阶段，表示获取锁出错，一般为redis连接问题
- false, true, nil/not nil 表示释放锁失败
- true, true, nil 表示执行成功