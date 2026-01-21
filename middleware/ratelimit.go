package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type limiter struct {
	windowStart time.Time
	count       int
	limit       int
	window      time.Duration
	mu          sync.Mutex
}

var (
	ipLimiters = make(map[string]*limiter) //map来制作局部限流器
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

func RedisRateLimiterMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()

		globalKey := "rate:global"
		ipKey := fmt.Sprintf("rate:ip:%s", ip)

		globalLimit := int64(15)
		ipLimit := int64(3)
		window := time.Minute

		globalCount, err := rdb.Incr(ctx, globalKey).Result()
		if err != nil {
			appErr := NewAppError(http.StatusInternalServerError, "RATE_LIMIT_STORAGE_ERROR", "rate limit storage error")
			handleError(c, appErr)
			c.Abort()
			return
		}
		if globalCount == 1 {
			rdb.Expire(ctx, globalKey, window)
		}
		ipCount, err := rdb.Incr(ctx, ipKey).Result()
		if err != nil {
			appErr := NewAppError(http.StatusInternalServerError, "RATE_LIMIT_STORAGE_ERROR", "rate limit storage error")
			handleError(c, appErr)
			c.Abort()
			return
		}
		if ipCount == 1 {
			rdb.Expire(ctx, ipKey, window)
		}

		if globalCount > globalLimit || ipCount > ipLimit {
			appErr := NewAppError(http.StatusTooManyRequests, "TOO_MANY_REQUEST", "too many request")
			handleError(c, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}
