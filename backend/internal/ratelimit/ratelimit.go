package ratelimit

import (
	"context"
	"time"

	"worldcup-mate/internal/database"

	"github.com/redis/go-redis/v9"
)

const (
	loginFailLimit = 5                // consecutive failures before freeze
	loginFailTTL   = 15 * time.Minute // freeze window
)

// AllowRequest enforces a fixed-window limit per key using Redis.
// Returns (allowed, error); on Redis failure it returns (false, err) so
// callers can degrade to allow (explicit degradation policy, SEC-02).
func AllowRequest(key string, limit int, window time.Duration) (bool, error) {
	ctx := context.Background()
	now := time.Now().Unix()
	countKey := "rl:" + key + ":count"
	windowKey := "rl:" + key + ":window"

	cur, err := database.RDB.Get(ctx, windowKey).Int64()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if cur != now {
		// New window: reset with TTL so keys never accumulate.
		if err := database.RDB.Set(ctx, windowKey, now, window).Err(); err != nil {
			return false, err
		}
		if err := database.RDB.Set(ctx, countKey, 1, window).Err(); err != nil {
			return false, err
		}
		return true, nil
	}
	count, err := database.RDB.Incr(ctx, countKey).Result()
	if err != nil {
		return false, err
	}
	return count <= int64(limit), nil
}

// LoginLocked reports whether the account is frozen after too many
// consecutive failed logins. Redis failures degrade to "not locked".
func LoginLocked(email string) (bool, error) {
	ctx := context.Background()
	count, err := database.RDB.Get(ctx, "rl:login:fail:"+email).Int()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return count >= loginFailLimit, nil
}

func RecordLoginFailure(email string) {
	ctx := context.Background()
	key := "rl:login:fail:" + email
	count, err := database.RDB.Incr(ctx, key).Result()
	if err != nil {
		return
	}
	if count == 1 {
		database.RDB.Expire(ctx, key, loginFailTTL)
	}
}

func ClearLoginFailures(email string) {
	database.RDB.Del(context.Background(), "rl:login:fail:"+email)
}
