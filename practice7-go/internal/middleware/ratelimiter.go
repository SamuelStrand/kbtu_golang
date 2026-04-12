package middleware

import (
	"net/http"
	"sync"
	"time"

	"practice7/internal/utils"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	Count     int
	ResetTime time.Time
}

type RateLimiter struct {
	limit  int
	window time.Duration
	secret []byte

	mu      sync.Mutex
	buckets map[string]*bucket
}

func NewRateLimiter(limit int, windowSeconds int, jwtSecret []byte) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  time.Duration(windowSeconds) * time.Second,
		secret:  jwtSecret,
		buckets: make(map[string]*bucket),
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rl.identifier(c)
		now := time.Now()

		rl.mu.Lock()
		b, ok := rl.buckets[id]
		if !ok {
			b = &bucket{Count: 0, ResetTime: now.Add(rl.window)}
			rl.buckets[id] = b
		}
		if now.After(b.ResetTime) {
			b.Count = 0
			b.ResetTime = now.Add(rl.window)
		}

		if b.Count >= rl.limit {
			resetAt := b.ResetTime
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":    "too many requests",
				"identity": id,
				"reset_at": resetAt.UTC().Format(time.RFC3339),
			})
			return
		}

		b.Count++
		rl.mu.Unlock()

		c.Next()
	}
}

func (rl *RateLimiter) identifier(c *gin.Context) string {
	if tokenStr := c.GetHeader("Authorization"); tokenStr != "" {
		claims, err := utils.ParseJWTLoose(tokenStr, rl.secret)
		if err == nil && claims.UserID != "" {
			return "user:" + claims.UserID
		}
	}
	return "ip:" + c.ClientIP()
}
