package ratelimiter

import (
	"net/http"

	"github.com/NetEase-Media/ngo/pkg/sentinel"
	"github.com/gin-gonic/gin"
)

func NewDefaultOptions() *RatelimiterOptions {
	return &RatelimiterOptions{
		DefaultMsg:    `{"code": -1, "msg":"http: Handler limit"}`,
		ErrorHttpCode: http.StatusTooManyRequests,
	}
}

func RateLimiter(opts ...Option) gin.HandlerFunc {
	o := NewDefaultOptions()

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(o)
	}

	return func(c *gin.Context) {
		e, b := sentinel.Entry(o.Resource)
		if b != nil {
			if o.ErrorHandler != nil {
				o.ErrorHandler(c)
				c.Abort()
			} else {
				c.AbortWithStatusJSON(o.ErrorHttpCode, o.DefaultMsg)
				return
			}
		} else {
			e.Exit()
			c.Next()
		}
	}
}
