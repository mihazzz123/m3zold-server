package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	store map[string][]time.Time
	mu    sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		store: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-1 * time.Hour)

		// Очистка старых записей
		var recentAttempts []time.Time
		for _, attempt := range rl.store[ip] {
			if attempt.After(windowStart) {
				recentAttempts = append(recentAttempts, attempt)
			}
		}

		// Максимум 5 попыток в час
		if len(recentAttempts) >= 100 {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}

		recentAttempts = append(recentAttempts, now)
		rl.store[ip] = recentAttempts

		c.Next()
	}
}

func getClientIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	return c.ClientIP()
}
