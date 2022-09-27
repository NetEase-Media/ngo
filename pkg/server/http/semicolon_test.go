package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func doSomething(c *gin.Context) {
	q, _ := c.GetQuery("q")
	log.Info(q)
	c.String(http.StatusOK, q)
}

func TestNoSemicolon(t *testing.T) {
	r := gin.New()
	r.GET("/", doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?q=a;b;c", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "a", w.Body.String())
}

func TestSemicolon(t *testing.T) {
	r := gin.New()
	r.Use(SemicolonMiddleware())
	r.GET("/", doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?q=a;b;c", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "a;b;c", w.Body.String())
}
