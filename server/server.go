// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/NetEase-Media/ngo/adapter/sentinel"

	"github.com/NetEase-Media/ngo/adapter/config"
	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/adapter/util"
	"github.com/NetEase-Media/ngo/adapter/xxljob"
	"github.com/NetEase-Media/ngo/client/db"
	"github.com/NetEase-Media/ngo/client/httplib"
	"github.com/NetEase-Media/ngo/client/kafka"
	"github.com/NetEase-Media/ngo/client/memcache"
	"github.com/NetEase-Media/ngo/client/multicache"
	"github.com/NetEase-Media/ngo/client/redis"
	"github.com/NetEase-Media/ngo/dlock"
	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
)

// 命令行参数
var (
	configPath string
	// 配置目录指定, 如果没指定 configPath，则会加载此目录下的app.yaml
	//configDir string
	ngofs = flag.NewFlagSet("ngoConfig", flag.ContinueOnError)
)

func initFlag() {
	ngofs.StringVar(&configPath, "c", "a", "config file path")
	//ngofs.StringVar(&configDir, "d", "", "config file directory")
}

type Server struct {
	*http.Server

	opt Options

	*gin.Engine

	// stopping 在开始停止server时被关闭，用来通知goroutine
	stopping chan struct{}
	stopped  chan struct{}

	// wgMutex 阻塞对wg的并发访问
	wgMutex sync.RWMutex
	// wg 用来等待goroutine完成后再关闭server
	wg sync.WaitGroup

	active      bool
	activeMutex sync.Mutex

	stopOnce sync.Once

	PreStart func() error
	PreStop  func(context.Context) error

	shutdownTimeout time.Duration
}

type MiddlewaresOptions struct {
	AccessLog *AccessLogMwOptions
}

type Options struct {
	Port            int
	Mode            string
	ShutdownTimeout time.Duration
	Middlewares     *MiddlewaresOptions
}

// PprofOptions 用于开启调试模式
type PprofOptions struct {
	Switch bool // 是否开启
	Port   int  // 配置端口号
}

func NewDefaultOptions() *Options {
	return &Options{
		Port:            8080,
		Mode:            gin.ReleaseMode,
		ShutdownTimeout: time.Second * 10,
		Middlewares: &MiddlewaresOptions{
			AccessLog: NewDefaultAccessLogOptions(),
		},
	}
}

type Method string

const (
	GET     Method = http.MethodGet
	HEAD    Method = http.MethodHead
	POST    Method = http.MethodPost
	PUT     Method = http.MethodPut
	PATCH   Method = http.MethodPatch
	DELETE  Method = http.MethodDelete
	CONNECT Method = http.MethodConnect
	OPTIONS Method = http.MethodOptions
	TRACE   Method = http.MethodTrace
)

var (
	server      *Server
	pprofServer *http.Server
)

func init() {
	initFlag()
}

// initConfig
func initConfig() {
	// config,error := conf.Init(getConfigFileName())
	err := config.Init(configPath)
	util.CheckError(err)
}

// initComponents 初始化所有外部组件
func initComponents() {
	// init service
	err := config.Unmarshal("service", &serviceOptions)
	util.CheckError(err)
	err = serviceOptions.Check()
	util.CheckError(err)

	// Init Logger
	// logOptions := log.NewDefaultOptions()
	var logOptions []log.Options
	err = config.Unmarshal("log", &logOptions)
	util.CheckError(err)
	// logOptions.FileName = serviceOptions.AppName
	err = log.Init(logOptions, serviceOptions.AppName)
	util.CheckError(err)

	// Init Redis
	var redisOptions []redis.Options
	err = config.Unmarshal("redis", &redisOptions)
	util.CheckError(err)
	err = redis.Init(redisOptions)
	util.CheckError(err)

	// Init Memcache
	var memOptions []memcache.Options
	err = config.Unmarshal("memcache", &memOptions)
	util.CheckError(err)
	err = memcache.Init(memOptions)
	util.CheckError(err)

	// Init DB
	var dbOptions []*db.Options
	err = config.Unmarshal("db", &dbOptions)
	util.CheckError(err)
	err = db.Init(dbOptions)
	util.CheckError(err)

	// Init Kafka
	var kafkaOptions = kafka.NewDefaultOptionsSlice(config.GetSliceSize("kafka"))
	err = config.Unmarshal("kafka", &kafkaOptions)
	util.CheckError(err)
	err = kafka.Init(kafkaOptions)
	util.CheckError(err)

	// Init Sentinel
	var sentinelOptions sentinel.Options
	err = config.Unmarshal("sentinel", &sentinelOptions)
	util.CheckError(err)
	err = sentinel.Init(&sentinelOptions)
	util.CheckError(err)

	// Init HTTPClient
	httplibOptions := httplib.NewDefaultOptions()
	err = config.Unmarshal("httpClient", httplibOptions)
	util.CheckError(err)
	httplib.Init(httplibOptions)

	// Init dlock
	var dlockOptions dlock.Options
	err = config.Unmarshal("dlock", &dlockOptions)
	util.CheckError(err)
	dlock.Init(dlockOptions)

	// Init pprof
	var pprof PprofOptions
	err = config.Unmarshal("pprof", &pprof)
	util.CheckError(err)
	if pprof.Switch { // 是否开启监听
		if pprof.Port == 0 {
			pprof.Port = 8899 // 默认调试端口号
		}
		pprofServer = &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", pprof.Port)}
		go func() {
			pprofServer.ListenAndServe()
		}()
	}

	// Init Miner
	var minerOptions []multicache.Options
	err = config.Unmarshal("multicache", &minerOptions)
	util.CheckError(err)
	multicache.Init(minerOptions)

	// Init xxljob
	var xxljobOptions xxljob.Options
	err = config.Unmarshal("xxljob", &xxljobOptions)
	util.CheckError(err)
	xxljob.Init(&xxljobOptions, serviceOptions.ClusterName)

}

// stopComponents 停止所有外部组件
func stopComponents(ctx context.Context) {
	// Stop Redis
	redis.StopAll()

	// Stop Kafka
	kafka.StopAll()

	// Close HTTPClient
	httplib.Close()

	// Stop pprof
	if pprofServer != nil {
		err := pprofServer.Shutdown(ctx)
		if err != nil {
			log.Errorf("shutting down pprof server error: %v", err)
		}
		log.Info("pprof server stopped")
	}

	// Stop xxljob
	xxljob.Stop()
}

func SetConfigPath(path string) {
	configPath = path
}

func GetConfigPath() string {
	return configPath
}

// Init 初始化全局唯一server
func Init() *Server {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()

	// 设置procs的数量为容器cpu的2倍
	runtime.GOMAXPROCS(runtime.GOMAXPROCS(0) * 2)

	ngofs.Parse(os.Args[1:])

	initConfig()

	initComponents()

	// Init ServerConfig
	opt := NewDefaultOptions()
	err := config.Unmarshal("httpServer", &opt)
	util.CheckError(err)

	server = newServer(opt)
	return server
}

func newServer(opt *Options) *Server {
	gin.SetMode(opt.Mode)
	engine := gin.New()

	if opt.Middlewares.AccessLog.Enabled {
		engine.Use(AccessLogMiddleware(opt.Middlewares.AccessLog))
	}

	engine.Use(OutermostRecover(),
		TrafficStopMiddleware(),
		ServerRecover(), SemicolonMiddleware())

	s := &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", opt.Port),
			Handler: engine,
		},
		Engine:          engine,
		opt:             *opt,
		stopping:        make(chan struct{}),
		stopped:         make(chan struct{}),
		shutdownTimeout: opt.ShutdownTimeout,
	}

	s.addServerHandler()
	return s
}

func (s *Server) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		os.Interrupt,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM, // kill -SIGTERM PID
	)

	go func() {
		if s.PreStart != nil {
			if err := s.PreStart(); err != nil {
				panic(fmt.Sprintf("start s failed: %s", err.Error()))
			}
		}
		log.WithField("port", s.opt.Port).Info("Start HTTPServer!")
		err := s.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {
			panic(fmt.Sprintf("start s failed: %s", err.Error()))
		}
	}()
	select {
	case <-ch:
		// 停止gin服务
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		s.Stop(ctx)
	case <-s.stopped:
		return
	}

}

// Stop 停止服务
func (s *Server) Stop(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()

	var err error
	log.Info("stopping server...")
	s.stopOnce.Do(func() {
		ch := make(chan struct{})

		go func() {
			if s.PreStop != nil {
				log.Info("PreStop start...")
				if e := s.PreStop(ctx); e != nil {
					log.Errorf("PreStop error: %v", e)
				}
				log.Info("PreStop end...")
			}

			s.stopServer(ctx)
			// 停止中间件
			stopComponents(ctx)
			close(ch)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-ch:
			close(s.stopped)
		}
	})
	log.Info("server stopped...")
	return err
}

// stopServer 停止web服务
func (s *Server) stopServer(ctx context.Context) {
	s.wgMutex.Lock()
	close(s.stopping)
	s.wgMutex.Unlock()
	s.wg.Wait() // 等待goroutine安全退出
	err := s.Shutdown(ctx)
	if err != nil {
		log.Errorf("shutting down http server error: %v", err)
	}
	log.Info("http server stopped")
}

// addServerHandler 注册服务状态相关route
func (s *Server) addServerHandler() *Server {
	health := s.Group("/health")
	health.GET("/online", s.onlineHandler)
	health.GET("/offline", s.offlineHandler)
	health.GET("/stop", s.offlineAndStopHandler)
	health.GET("/check", s.checkHandler)   // liveness probe
	health.GET("/status", s.statusHandler) // readiness probe
	return s
}

func (s *Server) AddRoute(method Method, path string, handlers ...gin.HandlerFunc) *Server {
	s.Handle(string(method), path, handlers...)
	return s
}

func (s *Server) AddRouteWithMethods(methods []Method, path string, handlers ...gin.HandlerFunc) *Server {
	if len(methods) == 0 {
		panic("methods can not be empty")
	}
	for i := range methods {
		s.Handle(string(methods[i]), path, handlers...)
	}
	return s
}

// StoppingNotify 返回一个channel，它会在server停止时被关闭
func (s *Server) StoppingNotify() <-chan struct{} { return s.stopping }

// GoAttach 用指定func创建一个goroutine，并用waitgroup跟踪它
// 注意传递的函数必须在s.StoppingNotify()关闭后停止
func (s *Server) GoAttach(f func()) {
	s.wgMutex.RLock() // 阻塞直至 close(s.stopping)
	defer s.wgMutex.RUnlock()
	select {
	case <-s.stopping:
		log.Warn("server has stopped; skipping GoAttach")
		return
	default:
	}

	// TODO: 这里可以增加一些监控代码
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

// offlineHandler 下线
func (s *Server) offlineHandler(c *gin.Context) {
	log.Info("[Offline] Server offline handler start...")
	s.activeMutex.Lock()
	s.active = false
	s.activeMutex.Unlock()
	if requestsFinished() {
		c.String(http.StatusOK, "ok")
		log.Info("[Offline] Requests finished...")
	} else {
		c.String(http.StatusBadRequest, "bad")
		log.Info("[Offline] Requests didnot finished...")
	}
	log.Info("[Offline] Server offline end...")
}

// offlineAndStopHandler 下线并且停服
func (s *Server) offlineAndStopHandler(c *gin.Context) {
	log.Info("[OfflineAndStop] Server offline and stop handler start...")
	s.activeMutex.Lock()
	s.active = false
	s.activeMutex.Unlock()
	if requestsFinished() {
		c.String(http.StatusOK, "ok")
		log.Info("[Offline] offline finished...")
		ctx, cancel := context.WithTimeout(c, s.shutdownTimeout)
		defer cancel()
		s.Stop(ctx)
	} else {
		c.String(http.StatusBadRequest, "bad")
		log.Info("[Offline] Requests didnot finished...")
	}
	log.Info("[OfflineAndStop] Server offline and stop end...")
}

func (s *Server) onlineHandler(c *gin.Context) {
	s.activeMutex.Lock()
	s.active = true
	s.activeMutex.Unlock()
	c.String(http.StatusOK, "ok")
	log.Info("Server online requested!")
}

func (s *Server) checkHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (s *Server) statusHandler(c *gin.Context) {
	s.activeMutex.Lock()
	active := s.active
	s.activeMutex.Unlock()
	if active {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusForbidden, "error")
	}
}
