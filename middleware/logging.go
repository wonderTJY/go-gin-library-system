package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		startAtstr := startAt.Format("2006-01-02 15:04:05")
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
				log.Print("request_count type assert to int64 failed")
			}
		} else {
			rc = 0
			log.Print("request_count haven't been set")
		}
		userID, ok := c.Get("user_id")
		if !ok {
			userID = 0
		}

		log.Printf("[%s] [rid:%s] method=%s path=%s status=%d latency=%dms ip=%s ua=%q rc=%d userID=%d",
			startAtstr, c.Writer.Header().Get(RequestIDHeader), method, path, statusCode, latency.Milliseconds(), clientIP, userAgent, rc, userID)

	}
}
