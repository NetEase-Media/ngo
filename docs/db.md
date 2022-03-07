# [Ngo](https://github.com/NetEase-Media/ngo)

---
## db
### 模块用途
提供数据库客户端，数据库客户端选择 [gorm](https://github.com/go-gorm/gorm) 实现。只需在配置中提供数据库配置，即可在运行中直接使用GetClient获取指定名字的客户端。其支持mysql、ddb两种模式的数据库连接，且都能自动上报哨兵监控数据。
### 使用说明
#### 配置
参数配置详见[db options](config.md#db-配置-dboptions)
#### 各模式应用
##### 获取实例
```go
// 配置文件方式
c = db.GetClient("db01")

// 手动实例化
opt := db.NewDefaultOptions()
opt.Name = "client1"
opt.Url = "xxxx"
c = db.NewClient(opt)
```
##### dsn 参数说明
```go
User             string            // Username
Passwd           string            // Password (requires User)
Net              string            // Network type
Addrs            []string          // Network address (requires Net)
DBName           string            // Database name
Params           map[string]string // Connection parameters
Collation        string            // Connection collation
Loc              *time.Location    // Location for time.Time values
MaxAllowedPacket int               // Max packet size allowed
ServerPubKey     string            // Server public key name
pubKey           *rsa.PublicKey    // Server public key
TLSConfig        string            // TLS configuration name
tls              *tls.Config       // TLS configuration
Timeout          time.Duration     // Dial timeout
ReadTimeout      time.Duration     // I/O read timeout
WriteTimeout     time.Duration     // I/O write timeout
Retry            int               // retry count with next qs, just for ddb

AllowAllFiles           bool // Allow all files to be used with LOAD DATA LOCAL INFILE
AllowCleartextPasswords bool // Allows the cleartext client side plugin
AllowNativePasswords    bool // Allows the native password authentication method
AllowOldPasswords       bool // Allows the old insecure password method
CheckConnLiveness       bool // Check connections for liveness before using them
ClientFoundRows         bool // Return number of matching rows instead of rows changed
ColumnsWithAlias        bool // Prepend table alias to column names
InterpolateParams       bool // Interpolate placeholders into query string
MultiStatements         bool // Allow multiple statements in one query
ParseTime               bool // Parse time values to time.Time
RejectReadOnly          bool // Reject read-only connections
```
*注意：如果是ddb，确定是否能使用服务端预处理，如果不能请设置InterpolateParams=true*

##### 执行命令
命令详见 [gorm 文档](https://gorm.io/zh_CN/docs/)
##### 关闭
```go
c.Close()
```

### 使用示例
- [examples/db](../examples/db) 