package middleware

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var totalRequests int64

func RequestCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		count := atomic.AddInt64(&totalRequests, 1)
		c.Set("request_count", count)
		c.Next()

	}
}
