//go:build !race
// +build !race

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func doPanic(ctx *gin.Context) {
	log.Panic("panic...")
}

func TestRecover(t *testing.T) {
	r := gin.New()
	r.Use(ServerRecover())

	r.GET("/", doPanic)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestRecover2(t *testing.T) {
	r := gin.New()
	r.Use(OutermostRecover())

	r.GET("/", doPanic)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", w.Body.String())
}
