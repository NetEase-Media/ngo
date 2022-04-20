package http

import (
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func TrafficStopMiddleware(requestCount *int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUri := c.Request.RequestURI
		defer func() {
			if !strings.HasPrefix(requestUri, "/health") {
				atomic.AddInt64(requestCount, -1)
			}
		}()

		if !strings.HasPrefix(requestUri, "/health") {
			atomic.AddInt64(requestCount, 1)
		}

		c.Next()
	}

}
