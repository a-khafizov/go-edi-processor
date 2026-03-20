package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-edi-document-processor/internal/logger"
	"go.uber.org/zap"
)

type Middleware struct {
	log *logger.Logger
}

func NewMiddleware(logger *logger.Logger) *Middleware {
	return &Middleware{log: logger}
}

func (m *Middleware) InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Only") != "true" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func (m *Middleware) Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

func (m *Middleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		m.log.Zap().Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
