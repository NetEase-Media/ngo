package http

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Options 是日志配置选项
type Options struct {
	Port            int
	Mode            string
	ShutdownTimeout time.Duration
	Middlewares     *MiddlewaresOptions
}

type MiddlewaresOptions struct {
	AccessLog  *AccessLogMwOptions
}

func NewDefaultOptions() *Options {
	return &Options{
		Port:            8080,
		Mode:            gin.DebugMode,
		ShutdownTimeout: 10 * time.Second,
		Middlewares: &MiddlewaresOptions{
			AccessLog:  NewDefaultAccessLogOptions(),
		},
	}
}

func checkOptions(opt *Options) error {
	return nil
}
