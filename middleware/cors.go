package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	origins := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		origins[o] = struct{}{}
	} //range取出，赋value给o，再把o填进map

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if len(origins) > 0 {
			if _, ok := origins[origin]; !ok {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
