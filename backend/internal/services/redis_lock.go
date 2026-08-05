package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"worldcup-mate/internal/database"
)

// REL-06: cross-instance distributed lock for sync tasks.
// - key carries provider/resource/competition/season
// - value is a unique owner id, so only the holder can release
// - TTL bounds the lock so a crashed holder never blocks forever
// - manual and scheduled syncs share the same key (same lock)

var (
	// ErrSyncAlreadyRunning marks "another worker holds the lock".
	ErrSyncAlreadyRunning = errors.New("sync already running")
)

// syncLockTTL bounds the lock so a crashed holder never blocks forever.
// It must exceed the longest possible sync run: the scheduled job's ctx
// timeout is 10min and a full round (matches + standings + 9s throttle)
// stays well below that — 20min keeps the lock alive across the entire
// run, so a late-finishing worker never overlaps a fresh one.
const syncLockTTL = 20 * time.Minute

// releaseLockScript deletes the key only when it still holds our owner
// value (atomic, so a stale holder can never release someone else's lock).
const releaseLockScript = `
if redis.call('get', KEYS[1]) == ARGV[1] then
	return redis.call('del', KEYS[1])
else
	return 0
end`

func syncLockKey(provider, resource, code string, season int) string {
	return "lock:sync:" + provider + ":" + resource + ":" + code + ":" + strconv.Itoa(season)
}

func randomLockOwner() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("150405.000000000")
	}
	return hex.EncodeToString(b)
}

// tryAcquireSyncLock attempts a non-blocking acquire. Returns (true, nil) on
// success, (false, nil) when another worker holds the lock, and (false, err)
// when Redis is unavailable (callers decide whether to degrade).
func tryAcquireSyncLock(ctx context.Context, key, owner string, ttl time.Duration) (bool, error) {
	if database.RDB == nil {
		// Defensive: the server always initializes RDB, but a nil client
		// must degrade like an outage instead of panicking (QA-03A).
		return false, fmt.Errorf("redis client is not initialized")
	}
	ok, err := database.RDB.SetNX(ctx, key, owner, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// releaseSyncLock releases the lock only if it still belongs to `owner`.
// Errors are logged, not propagated (a failed release self-heals via TTL).
func releaseSyncLock(ctx context.Context, key, owner string) {
	if _, err := database.RDB.Eval(ctx, releaseLockScript, []string{key}, owner).Result(); err != nil {
		log.Printf("[sync-lock] failed to release %s: %v", key, err)
	}
}
