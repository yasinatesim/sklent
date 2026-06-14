package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Rule struct {
	Burst    int
	Interval time.Duration
}

type bucket struct {
	tokens   int
	lastFill time.Time
}

// RateLimit is a per-IP token bucket. Demo-grade (in-memory); production would back it with Redis.
func RateLimit(rule Rule) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]*bucket{}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		b, ok := buckets[ip]
		now := time.Now()
		if !ok {
			b = &bucket{tokens: rule.Burst, lastFill: now}
			buckets[ip] = b
		}
		if now.Sub(b.lastFill) >= rule.Interval {
			b.tokens = rule.Burst
			b.lastFill = now
		}
		if b.tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate_limited"})
			return
		}
		b.tokens--
		mu.Unlock()
		c.Next()
	}
}
