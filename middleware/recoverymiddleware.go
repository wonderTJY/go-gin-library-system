package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic : (type:%T) %v \n%s", r, r, debug.Stack())
				appError := NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				handleError(c, appError)
				c.Abort()
			}
		}()
		c.Next()

	}
}
