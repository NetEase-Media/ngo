package admin

import (
	"context"

	"github.com/NetEase-Media/ngo/pkg/hooks"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/server/http"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func init() {
	hooks.Register(hooks.Init, NewAdmin)
}

var admin *http.Server

func Get() *http.Server {
	return admin
}

func NewAdmin(ctx context.Context) error {
	opt := http.NewDefaultOptions()
	opt.Port = 8899
	opt.Mode = ""
	opt.Middlewares.AccessLog.Enabled = true

	server, err := http.New(opt)
	if err != nil {
		return err
	}
	http.HealthCheck(server)
	pprof.Register(server.Engine)
	server.AddRoute(http.GET, "/logging", getLogger)
	server.AddRoute(http.POST, "/logging", postLogger)

	admin = server
	hooks.Register(hooks.ServerStart, func(ctx context.Context) error {
		return server.Start()
	})
	hooks.Register(hooks.ServerStop, func(ctx context.Context) error {
		return server.Stop()
	})
	hooks.Register(hooks.ServerGracefulStop, func(ctx context.Context) error {
		return server.GracefulStop(ctx)
	})
	return nil
}

func getLogger(c *gin.Context) {
	loggers := log.GetLoggers()
	r := make(map[string]string, len(loggers))
	for k, v := range loggers {
		r[k] = v.GetLevel().String()
	}
	if r[log.DefaultLoggerName] == "" {
		r[log.DefaultLoggerName] = log.DefaultLogger().GetLevel().String()
	}
	c.JSON(http.StatusOK, r)
}

func postLogger(c *gin.Context) {
	m := make(map[string]interface{})
	c.BindJSON(&m)
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if k == log.DefaultLoggerName {
			if level, err := log.ParseLevel(s); err == nil {
				log.SetLevel(level)
			}
		}
		if logger := log.GetLogger(k); logger != nil {
			if level, err := log.ParseLevel(s); err != nil {
				logger.SetLevel(level)
			}
		}
	}
	c.JSON(http.StatusOK, nil)
}
