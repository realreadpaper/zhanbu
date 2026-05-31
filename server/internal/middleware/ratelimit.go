package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// RateLimiter is a simple in-memory rate limiter using the token bucket pattern.
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]*bucket
	rate     int           // max requests per interval
	interval time.Duration // time window
}

type bucket struct {
	count    int
	resetAt  time.Time
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*bucket),
		rate:     rate,
		interval: interval,
	}
}

// Allow checks if a request from the given key is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.requests[key]

	if !exists || now.After(b.resetAt) {
		// New or expired bucket
		rl.requests[key] = &bucket{
			count:   1,
			resetAt: now.Add(rl.interval),
		}
		return true
	}

	if b.count >= rl.rate {
		return false
	}

	b.count++
	return true
}

// RateLimitMiddleware returns a Gin middleware that limits requests by IP.
func RateLimitMiddleware(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	// Periodic cleanup of expired buckets
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			limiter.mu.Lock()
			now := time.Now()
			for key, b := range limiter.requests {
				if now.After(b.resetAt) {
					delete(limiter.requests, key)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.Allow(key) {
			response.Error(c, http.StatusTooManyRequests, apperrors.ErrRateLimited,
				apperrors.GetMessage(apperrors.ErrRateLimited))
			c.Abort()
			return
		}
		c.Next()
	}
}
