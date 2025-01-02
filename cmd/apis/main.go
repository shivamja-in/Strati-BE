package main

import (
	"os"

	"github.com/gin-gonic/gin"
	ginMiddlewares "github.com/shivamja-in/strati-be/internal/middlewares/telemetry"
	"github.com/shivamja-in/strati-be/pkg/telemetry"
)

func main() {
	logger := telemetry.Telemetry()

	server := gin.New()
	server.Use(ginMiddlewares.GinLoggerMiddleware(logger), ginMiddlewares.GinRecoveryMiddleware(logger))

	PORT := os.Getenv("PORT")

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	logger.Info().Str("port", PORT).Msg("Starting app ....")

	if server.Run() != nil {
		logger.Fatal().Msg("Server failed to work")
	}
}
