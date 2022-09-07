package ratelimiter

import (
	"github.com/gin-gonic/gin"
)

type Option func(*RatelimiterOptions)

type RatelimiterOptions struct {
	Resource      string
	DefaultMsg    string
	ErrorHttpCode int
	ErrorHandler  gin.HandlerFunc
}

// Optional parameters
func WithResource(resource string) Option {
	return func(o *RatelimiterOptions) {
		o.Resource = resource
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(o *RatelimiterOptions) {
		o.ErrorHttpCode = code
	}
}

// Optional parameters
func WithDefaultMsg(s string) Option {
	return func(o *RatelimiterOptions) {
		o.DefaultMsg = s
	}
}

// Optional parameters
func WithErrorHandler(f gin.HandlerFunc) Option {
	return func(o *RatelimiterOptions) {
		o.ErrorHandler = f
	}
}
