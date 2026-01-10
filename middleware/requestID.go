package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get(RequestIDHeader)
		if rid == "" {
			rid = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set(RequestIDHeader, rid)

		c.Next()
	}
}
