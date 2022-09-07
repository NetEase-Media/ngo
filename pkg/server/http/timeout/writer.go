package timeout

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Writer is a writer with memory buffer
type Writer struct {
	gin.ResponseWriter
	body         *bytes.Buffer
	headers      http.Header
	mu           sync.Mutex
	timeout      bool
	wroteHeaders bool
	code         int
}

// NewWriter will return a timeout.Writer pointer
func NewWriter(buf *bytes.Buffer) *Writer {
	return &Writer{body: buf, headers: make(http.Header)}
}

// Write will write data to response body
func (w *Writer) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timeout || w.body == nil {
		return 0, nil
	}
	return w.body.Write(data)
}

// WriteHeader will write http status code
func (w *Writer) WriteHeader(code int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	checkWriteHeaderCode(code)
	if w.timeout || w.wroteHeaders {
		return
	}
	w.writeHeader(code)
}

func (w *Writer) writeHeader(code int) {
	w.wroteHeaders = true
	w.code = code
}

// Header will get response headers
func (w *Writer) Header() http.Header {
	return w.headers
}

// WriteString will write string to response body
func (w *Writer) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

// FreeBuffer will release buffer pointer
func (w *Writer) FreeBuffer() {
	w.body = nil
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid http status code: %d", code))
	}
}
