package gin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GinLoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info().
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}

func GinRecoveryMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().
					Interface("error", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Msg("Recovered from panic")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
