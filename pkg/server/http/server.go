package http

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/NetEase-Media/ngo/pkg/client/httplib"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/gin-gonic/gin"
)

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

var defaultHttpServer *Server

func Set(server *Server) {
	defaultHttpServer = server
}

func Get() *Server {
	return defaultHttpServer
}

func New(opt *Options) (*Server, error) {
	if opt.Mode != "" {
		gin.SetMode(opt.Mode)
	}
	engine := gin.New()

	s := &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", opt.Port),
			Handler: engine,
		},
		Engine: engine,
		Opt:    opt,
	}

	engine.Use(OutermostRecover())
	engine.Use(TrafficStopMiddleware(&s.handleCount))
	if opt.Middlewares.AccessLog.Enabled {
		engine.Use(AccessLogMiddleware(opt.Middlewares.AccessLog))
	}
	engine.Use(ServerTraceMiddleware())
	engine.Use(ServerRecover())
	engine.Use(SemicolonMiddleware())
	s.AddRoute(GET, "/healthz", func(c *gin.Context) {
		c.String(200, "health")
	})
	return s, nil
}

type Server struct {
	*gin.Engine
	*http.Server
	Opt *Options

	healthy     func() bool
	active      int32
	handleCount int64
}

func (s *Server) AddRoute(method Method, path string, handlers ...gin.HandlerFunc) *Server {
	s.Handle(string(method), path, handlers...)
	return s
}

func (s *Server) AddRouteWithMethods(methods []Method, path string, handlers ...gin.HandlerFunc) *Server {
	if len(methods) == 0 {
		log.Panic("methods can not be empty")
	}
	for i := range methods {
		s.Handle(string(methods[i]), path, handlers...)
	}
	return s
}

func (s *Server) Start() error {
	err := s.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Info("http server is closed")
		return nil
	}
	return err
}

func (s *Server) Stop() error {
	return s.Shutdown(context.Background())
}

func (s *Server) Healthz() bool {
	code, _ := httplib.Get(fmt.Sprintf("http://localhost:%d/healthz", s.Opt.Port)).Do(context.Background())
	return code == 200
}

func (s *Server) GracefulStop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.Opt.ShutdownTimeout)
	defer cancel()
	return s.Shutdown(ctx)
}

func (s *Server) requestsFinished() bool {
	return atomic.LoadInt64(&s.handleCount) == 0
}
