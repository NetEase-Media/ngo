//go:build !race
// +build !race

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTrafficStop(t *testing.T) {
	var handleCount int64

	r := gin.New()
	r.Use(TrafficStopMiddleware(&handleCount))

	r.GET("/", func(context *gin.Context) {
		time.Sleep(1000 * time.Millisecond)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	go r.ServeHTTP(w, req)

	time.Sleep(500 * time.Millisecond)
	assert.False(t, handleCount == 0)
	time.Sleep(800 * time.Millisecond)
	assert.True(t, handleCount == 0)
}
