package http

import (
	"net"
	"strconv"
	"strings"

	"github.com/NetEase-Media/ngo/pkg/tracing"
	"github.com/gin-gonic/gin"
)

func ServerTraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tracing.Enabled() && !strings.HasPrefix(c.Request.RequestURI, "/health") {
			span, ctx := tracing.StartSpanFromCarrier(c.Request.Context(), "server", c.Request.Header)
			tracing.SpanKind.Set(span, "server")
			tracing.SpanType.Set(span, "HTTP_SERVER")
			tracing.HttpServerRequestUrl.Set(span, c.FullPath())
			tracing.HttpServerRequestHost.Set(span, c.Request.Host)
			tracing.HttpServerRequestMethod.Set(span, c.Request.Method)
			if remoteIp, remotePortStr, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
				tracing.HttpServerPeerHost.Set(span, remoteIp)
				if remotePort, err := strconv.Atoi(remotePortStr); err == nil {
					tracing.HttpServerPeerPort.Set(span, uint16(remotePort))
				}
			}

			tracing.HttpServerRequestPath.Set(span, c.Request.URL.Path)
			tracing.HttpServerRequestSize.Set(span, c.Request.ContentLength)

			c.Request = c.Request.WithContext(ctx)
			defer func() {
				code := c.Writer.Status()
				tracing.HttpServerResponseStatus.Set(span, uint16(code))
				tracing.HttpServerResponseSize.Set(span, c.Writer.Size())
				span.Finish()
			}()
		}
		c.Next()
	}
}
