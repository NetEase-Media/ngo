//+build !race

package timeout

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func errorResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "timeout")
}

func doSomething(c *gin.Context) {
	q, _ := c.GetQuery("t")
	t, _ := strconv.ParseInt(q, 10, 64)
	time.Sleep(time.Millisecond * time.Duration(t))
	c.String(http.StatusOK, "success")
}

func TestTimeout(t *testing.T) {
	r := gin.New()
	r.GET("/", Timeout(WithTimeout(50*time.Millisecond), WithHandler(doSomething), WithErrorHandler(errorResponse)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?t=100", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)
	assert.Equal(t, "timeout", w.Body.String())
}

func TestWithoutTimeout(t *testing.T) {
	r := gin.New()
	r.GET("/", Timeout(WithTimeout(50*time.Millisecond), WithHandler(doSomething), WithErrorHandler(errorResponse)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?t=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
}
