package timeout

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CallBackFunc func(*http.Request)

// Option for timeout
type Option func(*Options)

// WithTimeout set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(t *Options) {
		if timeout > 0 {
			t.timeout = timeout
		}
	}
}

// WithHandler set handle process
func WithHandler(f gin.HandlerFunc) Option {
	return func(t *Options) {
		t.handler = f
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(t *Options) {
		t.errorHttpCode = code
	}
}

// Optional parameters
func WithDefaultMsg(s string) Option {
	return func(t *Options) {
		t.defaultMsg = s
	}
}

// Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *Options) {
		t.callBack = f
	}
}

// Optional parameters
func WithErrorHandler(f gin.HandlerFunc) Option {
	return func(t *Options) {
		t.errorHandler = f
	}
}

// Options struct
type Options struct {
	timeout       time.Duration
	handler       gin.HandlerFunc
	errorHttpCode int
	defaultMsg    string
	errorHandler  gin.HandlerFunc
	callBack      CallBackFunc
}
