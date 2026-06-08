package middleware

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

type memoryRateEntry struct {
	count     int64
	expiresAt time.Time
}

var memoryRateLimits = struct {
	sync.Mutex
	items map[string]memoryRateEntry
}{items: map[string]memoryRateEntry{}}

func RateLimit(name string, limit int64, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	return func(c *gin.Context) {
		key := rateLimitKey(name, c)
		allowed, err := allowRequest(c.Request.Context(), key, limit, window)
		if err != nil {
			allowed = allowRequestInMemory(key, limit, window)
		}
		if !allowed {
			utils.Error(c, 429, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}

func rateLimitKey(name string, c *gin.Context) string {
	if userID, ok := c.Get("user_id"); ok {
		return fmt.Sprintf("rate:%s:user:%v", name, userID)
	}
	ip := net.ParseIP(c.ClientIP())
	if ip == nil {
		return fmt.Sprintf("rate:%s:ip:unknown", name)
	}
	return fmt.Sprintf("rate:%s:ip:%s", name, ip.String())
}

func allowRequest(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	if database.RDB == nil {
		return false, fmt.Errorf("redis unavailable")
	}
	n, err := database.RDB.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if n == 1 {
		_ = database.RDB.Expire(ctx, key, window).Err()
	}
	return n <= limit, nil
}

func allowRequestInMemory(key string, limit int64, window time.Duration) bool {
	now := time.Now()
	memoryRateLimits.Lock()
	defer memoryRateLimits.Unlock()

	entry := memoryRateLimits.items[key]
	if now.After(entry.expiresAt) {
		entry = memoryRateEntry{expiresAt: now.Add(window)}
	}
	entry.count++
	memoryRateLimits.items[key] = entry

	for itemKey, item := range memoryRateLimits.items {
		if now.After(item.expiresAt) {
			delete(memoryRateLimits.items, itemKey)
		}
	}
	return entry.count <= limit
}
