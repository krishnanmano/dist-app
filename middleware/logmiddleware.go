package middleware

import (
	"dist-app/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := time.Since(start)

		fields := []logger.Field{
			logger.String("client_ip", c.ClientIP()),
			logger.String("latency", fmt.Sprintf("%v", duration)),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.RequestURI),
			logger.String("status", fmt.Sprintf("%d", c.Writer.Status())),
			//logger.String("user_id", c.Request.Header["user_id"][0]),
			logger.String("referrer", c.Request.Referer()),
			logger.String("request_id", c.Writer.Header().Get("Request-Id")),
			// "api_version": util.ApiVersion,
		}

		if c.Writer.Status() >= 500 {
			logger.Log.Error(c.Errors.String())
		} else {
			logger.Log.Info("asdasd", fields...)
		}
	}
}
