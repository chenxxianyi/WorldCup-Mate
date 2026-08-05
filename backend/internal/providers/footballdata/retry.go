package footballdata

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// REL-05: retry policy for the external football-data API.
// - retries transient failures: network/timeout errors, 429, 502, 503, 504
// - never retries client-side 4xx (parameter mistakes must fail fast)
// - exponential backoff with random jitter, honoring Retry-After on 429
// - context cancellation aborts immediately (no retry)

// RetryConfig controls the backoff schedule.
type RetryConfig struct {
	MaxRetries int           // total attempts = MaxRetries + 1
	BaseDelay  time.Duration // delay for the first retry
	MaxDelay   time.Duration // cap for exponential growth
}

func defaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   5 * time.Second,
	}
}

func retryableStatus(code int) bool {
	return code == http.StatusTooManyRequests ||
		code == http.StatusBadGateway ||
		code == http.StatusServiceUnavailable ||
		code == http.StatusGatewayTimeout
}

// retryDelay returns base*2^attempt (capped) with ±50% jitter so parallel
// sync workers do not stampede the upstream API in lockstep.
func retryDelay(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay >= maxDelay {
			delay = maxDelay
			break
		}
	}
	half := delay / 2
	if half <= 0 {
		return delay
	}
	return half + time.Duration(rand.Int63n(int64(half)))
}

// parseRetryAfter reads the Retry-After header: either seconds or an
// HTTP-date. Returns 0 when absent/unparseable (caller falls back to
// exponential backoff).
func parseRetryAfter(header string, now time.Time) time.Duration {
	if header == "" {
		return 0
	}
	if seconds, err := time.ParseDuration(header + "s"); err == nil && seconds > 0 {
		return seconds
	}
	if t, err := http.ParseTime(header); err == nil {
		if d := t.Sub(now); d > 0 {
			return d
		}
	}
	return 0
}

// do performs the request with the retry policy. The caller owns closing
// the returned body. `req` must be safe to reuse (GET, no body consumed).
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	cfg := c.retry
	var lastErr error
	var resp *http.Response

	for attempt := 0; ; attempt++ {
		resp = nil
		r, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				// Cancelled/deadline: never retry.
				return nil, ctx.Err()
			}
			lastErr = fmt.Errorf("football-data request failed: %w", err)
		} else {
			if r.StatusCode >= 200 && r.StatusCode < 300 {
				return r, nil
			}
			lastErr = fmt.Errorf("football-data returned status %d", r.StatusCode)
			resp = r
			r.Body.Close()
			if !retryableStatus(r.StatusCode) {
				// 4xx (except 429) and other non-retryable codes: fail fast.
				return nil, lastErr
			}
		}

		if attempt >= cfg.MaxRetries {
			return nil, lastErr
		}

		delay := retryDelay(attempt, cfg.BaseDelay, cfg.MaxDelay)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			if ra := parseRetryAfter(resp.Header.Get("Retry-After"), time.Now()); ra > 0 && ra < delay {
				delay = ra // honor the upstream's own wait hint
			}
		}
		if err := sleepContext(ctx, delay); err != nil {
			return nil, err
		}
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
