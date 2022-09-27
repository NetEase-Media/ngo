package http

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func SemicolonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" && strings.Contains(rawQuery, ";") {
			c.Request.URL.RawQuery = strings.ReplaceAll(rawQuery, ";", "%3B")
		}
		c.Next()
	}
}
