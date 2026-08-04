package middleware

import (
	"log"
	"time"

	"worldcup-mate/internal/ratelimit"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

// RateLimit enforces a fixed-window limit per client IP using Redis.
// When Redis is unavailable the request is allowed through (explicit
// degradation, logged) so auth never hard-fails on cache outages (SEC-02).
func RateLimit(keyPrefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Key is bound to the client IP; a constant prefix alone would turn
		// this into a global shared counter (blocking bug).
		key := keyPrefix + ":" + c.ClientIP()
		allowed, err := ratelimit.AllowRequest(key, limit, window)
		if err != nil {
			log.Printf("[rate-limit] redis unavailable, allowing request: %v", err)
			c.Next()
			return
		}
		if !allowed {
			utils.Error(c, 429, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
