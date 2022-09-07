//go:build !race
// +build !race

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NetEase-Media/ngo/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	metrics.Init("test", "test")

	r := gin.New()
	r.Use(UrlMetricsMiddleware(NewDefaultUrlMetricsOptions()))

	r.GET("/", func(context *gin.Context) {

	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
