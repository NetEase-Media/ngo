//go:build !race
// +build !race

package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func errorResponse(c *gin.Context) {
	c.String(http.StatusTooManyRequests, "limit exceeded")
}

func doSomething(c *gin.Context) {
	c.String(http.StatusOK, "success")
}

func TestLimit(t *testing.T) {
	r := gin.New()
	r.GET("/", RateLimiter(WithResource("abc"), WithErrorHandler(errorResponse)), doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Equal(t, "limit exceeded", w2.Body.String())
}

func TestGlobalMW(t *testing.T) {
	r := gin.New()
	r.Use(RateLimiter(WithResource("def"), WithErrorHandler(errorResponse)))
	r.GET("/", doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Equal(t, "limit exceeded", w2.Body.String())
}
