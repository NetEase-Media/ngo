
# [Ngo](https://github.com/NetEase-Media/ngo)

---

## cron 
* 主要用作单机定时任务
* 包装了“robfig/cron”类库

### cron 功能

#### 任务管理
```go
c := ngocron.New()
c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
c.AddFunc("CRON_TZ=Asia/Shanghai 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
c.AddFunc("@hourly",      func() { fmt.Println("Every hour, starting an hour from now") })
c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
c.Start()
// Start之后也可以添加任务
c.AddFunc("@daily", func() { fmt.Println("Every day") })
//停止任务调度，不会停止正在运行的任务
c.Stop()
```

### cron 表达式
#### cron表达式介绍
cron表达式一般由空格分隔的有5个字段

|字段名       | 是否必须   | 值范围          | 允许的特殊字符|
|----------   | ---------- | --------------  | ---------------|
|Minutes      | 是        | 0-59            | * / , -|
|Hours        | 是| 0-23            | * / , -|
|Day of month | 是| 1-31            | * / , - ?|
|Month        | 是| 1-12 or JAN-DEC | * / , -|
|Day of week  | 是| 0-6 or SUN-SAT  | * / , - ?|

其他可参考 
* [robfig/cron go doc](https://pkg.go.dev/github.com/robfig/cron/v3)
* [Cron Wiki](https://en.wikipedia.org/wiki/Cron)
#### 其他表达式格式
可自定义表达式格式，支持秒级
```go
cron.New(
	cron.WithParser(
		cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
```
内置秒级表达式支持，如果使用，需要按以下初始化
```go
cron.New(cron.WithSeconds())
```

#### 特殊字符
1. 星号（*）代表任何可能的值。例如：在“小时域”里的星号等于是“每一个小时”。
2. 逗号（,）表示列举，例如： 1,3,4,7 * * * * 表示每小时的1、3、4、7分会执行任务。
3. 连词符（-）表示范围，例如：1-6 * * * * 表示每小时的1到6分钟内，每分钟都会执行任务。
4. 问号(%) 用户DayOfMonth和DayOfWeek,指定日期和指定星期存在冲突，所以指定日期后，指定星期一般用？反之，如果指定了日期，指定星期一般用？。
5. 斜线(/) 表示时间间隔步长，比如分钟字段使用0/15表示为第0，15，30，45分钟

#### 预定义日程

| 值  | 描述 |备注|
| ------ | ------ |--|
|@yearly (or @annually) | 每年1月1号0点0分        | 0 0 1 1 *
|@monthly               | 每月1号0点0分 | 0 0 1 * *
|@weekly                | 每周日0点0分  | 0 0 * * 0
|@daily (or @midnight)  | 每天0点0分    | 0 0 * * *
|@hourly                | 每小时0分        | 0 * * * *
|@every <duration>| 每duration调用一次 [duration解释](https://golang.org/pkg/time/#ParseDuration) |

#### 时区
默认使用机器的时区，time.Local,也可以自定义时区
```go
cron.New(cron.WithLocation(time.UTC))
```
### Job Wrapper
job wrapper，相等于连接器
如果定义wrapper如下
NewChain(m1,m2,m3).Then(job)
相当于
m1(m2(m3(job)))
目前内置的job wrapper有
* Recover  记录job的panic并打印日志
* DelayIfStillRunning 如果上一个任务未执行完成，则延迟当前任务的执行，如果延迟超过1分支，则打印日志
* SkipIfStillRunning 如果上一个任务未执行完成，则跳过当前的任务的执行。
使用方式
```
//为所有job安装
cron.New(cron.WithChain(
	cron.SkipIfStillRunning(logger),
))
//为单独job安装
job = cron.NewChain(
	cron.SkipIfStillRunning(logger),
).Then(job)
```


### 其他功能
* 分布式环境只执行一次？建议使用分布式锁

## 代码设计
### 框架选型
使用框架 [robfig/cron](https://github.com/robfig/cron)
对原有的API进行了简化
![ngo定时任务代码设计](/uploads/cce7dada99f00fa10df5bd655318827f/ngo定时任务代码设计.png)
