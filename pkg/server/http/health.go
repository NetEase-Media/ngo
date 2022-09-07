package http

import (
	"net/http"
	"sync/atomic"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/gin-gonic/gin"
)

func HealthCheck(s *Server) {
	health := s.Group("/health")
	health.GET("/online", s.onlineHandler)
	health.GET("/offline", s.offlineHandler)
	health.GET("/check", s.checkHandler)   // liveness probe
	health.GET("/status", s.statusHandler) // readiness prob
}

func (s *Server) SetHealthyFn(fn func() bool) {
	s.healthy = fn
}

// offlineHandler 下线
func (s *Server) offlineHandler(c *gin.Context) {
	atomic.StoreInt32(&s.active, 0)
	if s.requestsFinished() {
		c.String(http.StatusOK, "ok")
		log.Info("Server offline requested!")
	} else {
		c.String(http.StatusBadRequest, "bad")
		log.Info("Server offline failed!")
	}
}

func (s *Server) onlineHandler(c *gin.Context) {
	atomic.StoreInt32(&s.active, 1)
	c.String(http.StatusOK, "ok")
	log.Info("Server online requested!")
}

func (s *Server) checkHandler(c *gin.Context) {
	ss := Get()
	if !ss.Healthz() {
		c.String(http.StatusForbidden, "error")
		return
	}

	if s.healthy != nil && !s.healthy() {
		c.String(http.StatusForbidden, "error")
		return
	}
	c.String(http.StatusOK, "ok")
	log.Info("Server check requested!")
}

func (s *Server) statusHandler(c *gin.Context) {
	if atomic.LoadInt32(&s.active) == 1 {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusForbidden, "error")
	}
}
