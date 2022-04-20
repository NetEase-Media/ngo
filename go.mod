module github.com/NetEase-Media/ngo

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/Shopify/sarama v1.27.2
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/agiledragon/gomonkey v2.0.2+incompatible
	github.com/alibaba/sentinel-golang v1.0.4
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/bluele/gcache v0.0.2
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/djimenez/iconv-go v0.0.0-20160305225143-8960e66bd3da
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-redis/redis/extra/rediscmd v0.2.0
	github.com/go-redis/redis/v8 v8.8.2
	github.com/go-zookeeper/zk v1.0.2
	github.com/golang/mock v1.6.0
	github.com/gomodule/redigo v1.8.8 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/mapstructure v1.4.3
	github.com/modern-go/reflect2 v1.0.2
	github.com/philchia/agollo/v4 v4.1.3
	github.com/pinpoint-apm/pinpoint-go-agent v0.5.1
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/uber/jaeger-lib v2.4.0+incompatible
	github.com/urfave/cli/v2 v2.4.1
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fasthttp v1.34.0
	github.com/xxl-job/xxl-job-executor-go v1.0.0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.etcd.io/etcd/client/v3 v3.5.4
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/time v0.0.0-20220411224347-583f2d630306 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.4
)

replace (
	github.com/pinpoint-apm/pinpoint-go-agent => github.com/machine3/pinpoint-go-agent v0.5.1-r1
	github.com/spf13/viper => github.com/machine3/viper v1.11.0-r2
)
