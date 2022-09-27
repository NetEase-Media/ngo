# [Ngo](https://github.com/NetEase-Media/ngo)

---
## Log
### 模块用途
用来记录项目中的日志，支持多级别、多文件
### 使用说明
日志使用logrus实现日志接口，并提供以下功能：
- 统一简洁的定制化格式输出，包含日志的时间、级别、代码行数、函数、日志体
- 可选按txt或json格式输出日志
- access、info、error日志分离到不同文件中
- 提供文件轮转功能，在日志文件达到指定大小或寿命后切换到新文件

服务默认输出txt的日志格式，样式如：

`时间 [级别] [代码目录/代码文件:行数] [函数名] [字段键值对] 日志体`

时间格式类似`2021-01-14 10:39:33.349`。

级别包含以下几种：
- panic
- fatal
- error
- warning
- info
- debug

如果未设置级别，被被默认设置为info。非测试状态不要开启debug，避免日志过多影响性能。
另外在日志输出时可以使用WithField或WithFields来字段的key-value，在创建子日志对象时可以用来清晰地辨认日志的使用范围，但平时尽量不要使用。另外如果要输出error也尽量避免使用字段，直接使用Error()方法输出为字符串是最快的。

#### 自定义日志输出

##### 自定义日志输出改进
从之前的支持单个输出到支持多个输出，配置的选项都不变。

#### 配置文件
* 新增format blank表示空白，不添加其他前缀或后缀
* 如果不写name,会被当作默认的日志输出
```yaml
log:
  - name: default
    path: ./log
    level: debug
    format: txt
    errorPath: ./error
  - name: haha
    path: ./log
    fileName: haha.log
    level: info
    format: blank
    errorPath: ./error
```
#### 代码中使用
```go
//如果使用默认logger,和之前类似
import "github.com/NetEase-Media/ngo/pkg/log"
log.Infof("123")
//如果使用自定义logger,需要以下面这种方式获取
import "github.com/NetEase-Media/ngo/pkg/log"

log.GetLogger("haha").Infof("123.....")
```

