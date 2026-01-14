package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	if !c.Writer.Written() {
		if appErr, ok := err.(*AppError); ok {
			c.JSON(appErr.StatusCode, gin.H{
				"code":       appErr.Code,
				"message":    appErr.Message,
				"request_id": c.GetString("request_id"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":       "INTERNAL_ERROR",
			"message":    "internal server error",
			"request_id": c.GetString("request_id"),
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
