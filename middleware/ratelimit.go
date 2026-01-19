package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type limiter struct {
	windowStart time.Time
	count       int
	limit       int
	window      time.Duration
	mu          sync.Mutex
}

var limter = &limiter{
	windowStart: time.Now(),
	limit:       3,
	window:      time.Minute,
}

func (l *limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	timenow := time.Now()
	if timenow.After(l.windowStart.Add(l.window)) {
		l.windowStart = time.Now()
		l.count = 0
	}
	if l.count < l.limit {
		l.count++
		return true
	}
	return false
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !limter.Allow() {
			err := NewAppError(http.StatusTooManyRequests, "TOO_MANY_REQUEST", "too many request")
			handleError(c, err)
			c.Abort()
			return
		}
		c.Next()

	}
}
