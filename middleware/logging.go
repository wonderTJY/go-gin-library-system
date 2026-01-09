package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		//startAtstr := time.Now().Format("2006-01-02 15:04:05")
		method := c.Request.Method
		path := c.Request.URL.Path
		header := c.Request.Header

		c.Next()

		statusCode := c.Writer.Status()
		latency := time.Since(startAt)
		log.Printf("method=%s path=%s header=%s status=%d latency=%v",
			method, path, header, statusCode, latency.Milliseconds())

	}
}
