package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/gin-gonic/gin"
)

var bufPool *BufferPool

const (
	defaultTimeout = 5 * time.Second
)

var defaultOptions = NewDefaultOptions()

func NewDefaultOptions() Options {
	return Options{
		callBack:      nil,
		defaultMsg:    `{"code": -1, "msg":"http: Handler timeout"}`,
		timeout:       defaultTimeout,
		errorHttpCode: http.StatusServiceUnavailable,
	}
}

// Timeout
func Timeout(opts ...Option) gin.HandlerFunc {
	t := defaultOptions

	// Loop through each option
	for _, opt := range opts {
		if opt == nil {
			panic("timeout Option not be nil")
		}
		// Call the option giving the instantiated
		opt(&t)
	}

	bufPool = &BufferPool{}

	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		cp := c.Copy()

		w := c.Writer
		buffer := bufPool.Get()
		tw := NewWriter(buffer)
		cp.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), t.timeout)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					log.Errorf("error: %v", p)
					panicChan <- p
				}
			}()
			t.handler(cp)
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.FreeBuffer()
			bufPool.Put(buffer)
			panic(p)

		case <-finish:
			c.Next()
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := w.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			w.WriteHeader(tw.code)
			if _, err := w.Write(buffer.Bytes()); err != nil {
				panic(err)
			}
			tw.FreeBuffer()
			bufPool.Put(buffer)
		case <-ctx.Done():
			c.Abort()
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.timeout = true
			tw.FreeBuffer()
			bufPool.Put(buffer)
			if t.errorHandler != nil {
				t.errorHandler(c)
			} else {
				c.String(t.errorHttpCode, t.defaultMsg)
			}
			if t.callBack != nil {
				t.callBack(c.Request.Clone(context.Background()))
			}
		}
	}
}
