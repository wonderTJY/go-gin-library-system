package middleware

import (
	"net/http"
	"runtime/debug"
	"trae-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L.Error("Panic Recovered",
					zap.Any("error", r),
					zap.String("stack", string(debug.Stack())), // 记录堆栈
					zap.String("rid", c.GetString("request_id")),
				)
				appError := NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				handleError(c, appError)
				c.Abort()
			}
		}()
		c.Next()

	}
}
