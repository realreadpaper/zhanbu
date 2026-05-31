package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// LoggerMiddleware returns a Gin middleware that logs requests using zerolog.
func LoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		event := logger.Info()
		if statusCode >= 400 {
			event = logger.Warn()
		}
		if statusCode >= 500 {
			event = logger.Error()
		}

		event.
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Int("body_size", c.Writer.Size()).
			Msg("request completed")
	}
}
