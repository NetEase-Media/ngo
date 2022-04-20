package http

//
//import (
//	"g.hz.netease.com/agent/nss-go-agent/collector/collectors/url"
//	"github.com/NetEase-Media/ngo/pkg/metrics"
//	collectors "github.com/NetEase-Media/ngo/pkg/metrics/colloctors"
//	"github.com/gin-gonic/gin"
//)
//
//type UrlMetricsMwOptions struct {
//	Enabled      bool
//	OriginalPath bool
//}
//
//func NewDefaultUrlMetricsOptions() *UrlMetricsMwOptions {
//	return &UrlMetricsMwOptions{
//		Enabled:      true,
//		OriginalPath: false,
//	}
//}
//
//func UrlMetricsMiddleware(opt *UrlMetricsMwOptions) gin.HandlerFunc {
//	if opt == nil {
//		opt = NewDefaultUrlMetricsOptions()
//	}
//	return func(c *gin.Context) {
//		if !opt.Enabled {
//			c.Next()
//			return
//		}
//		if !metrics.IsMetricsEnabled() {
//			c.Next()
//			return
//		}
//
//		var path string
//		// 是否取原始地址
//		if opt.OriginalPath {
//			path = c.Request.RequestURI
//		} else {
//			path = c.FullPath()
//			if len(path) == 0 {
//				path = c.Request.RequestURI
//			}
//		}
//
//		requestMethod := c.Request.Method
//
//		var (
//			stats      *url.StatsHolder
//			statusCode = 200
//		)
//		defer func() {
//			if err := recover(); err != nil {
//				if e, ok := err.(error); ok {
//					collectors.UrlCollector().OnError(stats, e)
//				}
//			}
//			collectors.UrlCollector().OnComplete(stats, statusCode)
//		}()
//
//		stats = collectors.UrlCollector().OnStart(path, requestMethod)
//		c.Next()
//		statusCode = c.Writer.Status()
//	}
//}
