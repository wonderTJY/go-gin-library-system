package middleware

import (
	"net/http"
	"trae-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppError struct {
	StatusCode int
	Code       string
	Message    string
}

// 实现error接口，让AppError成为error
func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statuscode int, code string, msg string) *AppError {
	return &AppError{
		StatusCode: statuscode,
		Code:       code,
		Message:    msg,
	}
}

func handleError(c *gin.Context, err error) {
	rid := c.GetString("request_id")

	if appErr, ok := err.(*AppError); ok {
		// 日志照常记录
		logger.L.Warn("Business Error",
			zap.String("rid", rid),
			zap.String("code", appErr.Code),
			zap.String("msg", appErr.Message))

		// 响应只在没写过的时候写
		if !c.Writer.Written() {
			c.JSON(appErr.StatusCode, gin.H{
				"code":       appErr.Code,
				"message":    appErr.Message,
				"request_id": rid,
			})
		}
		return
	}

	// 日志照常记录
	logger.L.Error("Internal Server Error",
		zap.String("rid", rid),
		zap.Error(err),
	)

	// 响应只在没写过的时候写
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":       "INTERNAL_ERROR",
			"message":    "internal server error",
			"request_id": rid,
		})
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err
		handleError(c, err)
	}
}
