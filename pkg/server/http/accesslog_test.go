package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAccessLog(t *testing.T) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
	r.GET("/ping", func(c *gin.Context) {
		c.Set("data", "data...")
		time.Sleep(3 * time.Millisecond)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	w := util.PerformRequest(r, "GET", "/ping?p=aaaa&q=bbbb", util.Header{Key: "X-Real-IP", Value: "1.1.1.1"},
		util.Header{Key: "User-Agent", Value: "AHC/2.1 ..."})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccessLog_branchCov(t *testing.T) {
	r := gin.Default()
	r.Use(AccessLogMiddleware(nil))
	opt := AccessLogMwOptions{
		Enabled: false,
		Pattern: `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`,
	}
	r.Use(AccessLogMiddleware(&opt))
	opt = AccessLogMwOptions{
		Enabled: true,
		Pattern: `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`,
		NoFile:  false,
	}
	r.Use(AccessLogMiddleware(&opt))
}

func TestAccessLog_(t *testing.T) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
}

func BenchmarkAccessLog(b *testing.B) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = "common"
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	util.RunRequest(b, r, "GET", "/ping?p=aaaa", util.Header{Key: "X-Real-IP", Value: "1.1.1.1"},
		util.Header{Key: "User-Agent", Value: "AHC/2.1"})
}
