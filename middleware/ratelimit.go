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

var (
	ipLimiters = make(map[string]*limiter)
	ipMu       sync.RWMutex
)
var globleLimter = &limiter{
	windowStart: time.Now(),
	limit:       15,
	window:      time.Minute,
}

func getLimiterForIP(ip string) *limiter {
	ipMu.RLock()
	l, ok := ipLimiters[ip]
	ipMu.RUnlock()
	if ok {
		return l
	}

	ipMu.Lock()
	defer ipMu.Unlock()

	if l, ok = ipLimiters[ip]; ok {
		return l
	}
	l = &limiter{
		windowStart: time.Now(),
		limit:       3,
		window:      time.Minute,
	}
	ipLimiters[ip] = l
	return l

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
		ip := c.ClientIP()
		limiter := getLimiterForIP(ip)

		if !globleLimter.Allow() || !limiter.Allow() {
			err := NewAppError(http.StatusTooManyRequests, "TOO_MANY_REQUEST", "too many request")
			handleError(c, err)
			c.Abort()
			return
		}
		c.Next()

	}
}
