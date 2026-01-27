package middleware

import (
	"time"
	"trae-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		statusCode := c.Writer.Status()
		latency := time.Since(startAt)
		rcTemp1, exist := c.Get("request_count")
		var rc int64
		if exist {
			if rcTemp2, ok := rcTemp1.(int64); ok {
				rc = rcTemp2
			} else {
				rc = 0
				logger.L.Warn("request_count type assert to int64 failed")
			}
		} else {
			rc = 0
			logger.L.Debug("request_count haven't been set")
		}
		userID, ok := c.Get("user_id")
		if !ok {
			userID = 0
		}

		logger.L.Info("HTTP Request",
			zap.String("rid", c.Writer.Header().Get(RequestIDHeader)),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
			zap.String("ua", userAgent),
			zap.Int64("rc", rc),
			zap.Any("userID", userID),
		)

	}
}
