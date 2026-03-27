package middlewares

import (
	"log/slog"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
)

func GenAndAttachRequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = utils.GenUUID()
		}
		c.Header("X-Request-ID", id)
		c.Set("requestId", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

func AttachLoggerToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.GetString("requestId")
		reqLogger := slog.Default().With(
			slog.String("request_id", reqId),
		)
		c.Set("logger", reqLogger)
		c.Next()
	}
}

func PerformBasicCalculations() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		logger := utils.GetLogger(c)

		logger.Info("request started", "method", c.Request.Method, "path", path)
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		msg := "request completed"
		attrs := []any{
			"status", status,
			"latency", latency,
			"path", path,
		}
		if status >= 500 {
			logger.Error(msg, attrs...)
		} else if status >= 400 {
			logger.Warn(msg, attrs...)
		} else {
			logger.Info(msg, attrs...)
		}
	}
}
