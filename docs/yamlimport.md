# [Ngo](https://github.com/NetEase-Media/ngo)

---

## 多环境配置&&配置导入
### 模块用途
该功能用来支持多个环境下定义文件目录，以及yaml文件拆分、导入

### 使用说明
#### 多环境配置目录
* 新增-d启动项可以设置配置文件目录，默认值为 (当前运行的go程序的路径/config)
* 比如 -d /path/to/directory/test 代表配置文件目录在/path/to/directory/test,  
  环境为test，会加载/path/to/directory/test 目录下配置文件，
  默认如果不指定-c参数，则会自动加载-d参数下配置文件： /path/to/directory/test/app.yaml


**！！！注意，如果-c, -d全都不加，那么默认加载的配置文件路径为：
当前运行的go程序的路径/config/app.yaml**

#### 加载配置文件目录下的其他配置文件
```go
// 导包
import "github.com/NetEase-Media/ngo/pkg/config"


// 读取配置目录中viper支持的文件，例如:json,yaml,yml,toml,properties,props,prop,hcl,ini等
// 返回值为(*Config, error) 
c, err := config.New("data.properties")

```

#### 配置文件导入功能
* 可以在yaml类型的配置文件中增加顶层配置configImports，用于包含其他yaml类型文件的的配置
  比如以下示例，会额外加载kafka.yaml 和 db.yaml
```yaml
//app.yaml
service:
  serviceName: app-test-service
  appName: app-test
  clusterName: app-test-cm_test
configImports:
  - data.properties
  - db.yaml
  - kafka.yaml

// db.yaml
db1:
  host: localhost
  port: 9090

// kafka.yaml
kafka1:
  version: 100 
```
可以再代码中用以下方法获取配置
```go
//导包
import "github.com/NetEase-Media/ngo/pkg/config"

//获取为自定义的struct
config.Unmarshal("kafka1", kafkaStructPointer)
```