package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"worldcup-mate/internal/database"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// withMiniredis swaps the global RDB for an in-memory Redis for the test.
func withMiniredis(t *testing.T, fn func(mr *miniredis.Miniredis)) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	old := database.RDB
	database.RDB = client
	defer func() {
		database.RDB = old
		_ = client.Close()
	}()
	fn(mr)
}

func TestSyncLockMutualExclusion(t *testing.T) {
	withMiniredis(t, func(mr *miniredis.Miniredis) {
		key := syncLockKey("football-data", "matches:PL", "PL", 2025)
		ctx := context.Background()

		ok, err := tryAcquireSyncLock(ctx, key, "owner-a", syncLockTTL)
		if err != nil || !ok {
			t.Fatalf("first acquire: ok=%v err=%v", ok, err)
		}
		ok, err = tryAcquireSyncLock(ctx, key, "owner-b", syncLockTTL)
		if err != nil || ok {
			t.Fatalf("second acquire must fail while held: ok=%v err=%v", ok, err)
		}

		releaseSyncLock(ctx, key, "owner-a")
		ok, err = tryAcquireSyncLock(ctx, key, "owner-b", syncLockTTL)
		if err != nil || !ok {
			t.Fatalf("acquire after release: ok=%v err=%v", ok, err)
		}
	})
}

func TestSyncLockOnlyHolderCanRelease(t *testing.T) {
	withMiniredis(t, func(mr *miniredis.Miniredis) {
		key := syncLockKey("football-data", "matches:PL", "PL", 2025)
		ctx := context.Background()

		if ok, _ := tryAcquireSyncLock(ctx, key, "owner-a", syncLockTTL); !ok {
			t.Fatal("acquire failed")
		}
		// A stale holder (owner-b) must NOT be able to release the lock.
		releaseSyncLock(ctx, key, "owner-b")
		ok, err := tryAcquireSyncLock(ctx, key, "owner-c", syncLockTTL)
		if err != nil || ok {
			t.Fatalf("lock was released by a non-holder: ok=%v err=%v", ok, err)
		}
		// The real holder still can.
		releaseSyncLock(ctx, key, "owner-a")
		ok, _ = tryAcquireSyncLock(ctx, key, "owner-c", syncLockTTL)
		if !ok {
			t.Fatal("lock not released by its holder")
		}
	})
}

func TestSyncLockExpiresByTTL(t *testing.T) {
	withMiniredis(t, func(mr *miniredis.Miniredis) {
		key := syncLockKey("football-data", "matches:PL", "PL", 2025)
		ctx := context.Background()

		if ok, _ := tryAcquireSyncLock(ctx, key, "crashed-worker", syncLockTTL); !ok {
			t.Fatal("acquire failed")
		}
		// A crashed holder never releases; the TTL must self-heal.
		mr.FastForward(syncLockTTL + time.Second)
		ok, err := tryAcquireSyncLock(ctx, key, "new-worker", syncLockTTL)
		if err != nil || !ok {
			t.Fatalf("lock did not expire: ok=%v err=%v", ok, err)
		}
	})
}

func TestSyncLockKeysAreIndependent(t *testing.T) {
	withMiniredis(t, func(mr *miniredis.Miniredis) {
		ctx := context.Background()
		plKey := syncLockKey("football-data", "matches:PL", "PL", 2025)
		saKey := syncLockKey("football-data", "matches:SA", "SA", 2025)

		if ok, _ := tryAcquireSyncLock(ctx, plKey, "a", syncLockTTL); !ok {
			t.Fatal("PL acquire failed")
		}
		ok, err := tryAcquireSyncLock(ctx, saKey, "b", syncLockTTL)
		if err != nil || !ok {
			t.Fatalf("different league must not be blocked: ok=%v err=%v", ok, err)
		}
		// Same league, different season -> different lock.
		plNextKey := syncLockKey("football-data", "matches:PL", "PL", 2026)
		ok, _ = tryAcquireSyncLock(ctx, plNextKey, "c", syncLockTTL)
		if !ok {
			t.Fatal("next season must have its own lock")
		}
	})
}

func TestSyncLockKeyFormat(t *testing.T) {
	key := syncLockKey("football-data", "matches:PL", "PL", 2025)
	want := "lock:sync:football-data:matches:PL:PL:2025"
	if key != want {
		t.Errorf("key = %q, want %q", key, want)
	}
	if !strings.Contains(key, "PL") || !strings.Contains(key, "2025") {
		t.Errorf("key must carry provider/resource/code/season: %q", key)
	}
}
