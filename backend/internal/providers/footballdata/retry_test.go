package footballdata

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// countServer wraps a handler that counts requests.
func countServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *int32) {
	t.Helper()
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		handler(w, r)
	}))
	return server, &hits
}

func TestDoRetriesTransientFailures(t *testing.T) {
	attempts := 0
	server, hits := countServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusServiceUnavailable) // 503
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	client.retry = RetryConfig{MaxRetries: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	resp, err := client.do(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if *hits != 2 {
		t.Errorf("hits = %d, want 2 (one 503 + one success)", *hits)
	}
}

func TestDoHonorsRetryAfter(t *testing.T) {
	attempts := 0
	server, hits := countServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests) // 429
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	// First backoff spans [2s, 4s] — strictly above Retry-After=1s, so the
	// assertion below can only pass if the Retry-After hint actually won.
	client.retry = RetryConfig{MaxRetries: 3, BaseDelay: 4 * time.Second, MaxDelay: 8 * time.Second}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	start := time.Now()
	resp, err := client.do(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success after 429+Retry-After, got %v", err)
	}
	defer resp.Body.Close()
	if *hits != 2 {
		t.Errorf("hits = %d, want 2", *hits)
	}
	if elapsed := time.Since(start); elapsed < 900*time.Millisecond || elapsed > 1500*time.Millisecond {
		t.Errorf("Retry-After=1s not honored: elapsed %v (want ~1s, not the 2-4s backoff)", elapsed)
	}
}

func TestDoExhaustsRetries(t *testing.T) {
	server, hits := countServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway) // 502 always
	})
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	client.retry = RetryConfig{MaxRetries: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	_, err := client.do(context.Background(), req)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if !strings.Contains(err.Error(), "502") {
		t.Errorf("error %q should mention the status", err)
	}
	if *hits != 4 {
		t.Errorf("hits = %d, want 4 (1 attempt + 3 retries)", *hits)
	}
}

func TestDoDoesNotRetryClient4xx(t *testing.T) {
	server, hits := countServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest) // 400: parameter error
	})
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	client.retry = RetryConfig{MaxRetries: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	_, err := client.do(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for 400")
	}
	if *hits != 1 {
		t.Errorf("hits = %d, want 1 (4xx must not be retried)", *hits)
	}
}

func TestDoStopsOnContextCancel(t *testing.T) {
	server, hits := countServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	client.retry = RetryConfig{MaxRetries: 5, BaseDelay: time.Hour, MaxDelay: time.Hour}

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err := client.do(ctx, req)
	if err == nil {
		t.Fatal("expected error after cancel")
	}
	if *hits != 1 {
		t.Errorf("hits = %d, want 1 (cancel must abort before the next retry)", *hits)
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		header string
		want   time.Duration
	}{
		{"", 0},
		{"5", 5 * time.Second},
		{"Mon, 02 Jan 2006 15:04:05 GMT", 0}, // past HTTP-date -> 0
	}
	// A future HTTP-date maps to the exact remaining duration.
	future := now.Add(90 * time.Second).Format(http.TimeFormat)
	if got := parseRetryAfter(future, now); got != 90*time.Second {
		t.Errorf("future http-date: got %v, want 90s", got)
	}
	for _, c := range cases {
		if got := parseRetryAfter(c.header, now); got != c.want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", c.header, got, c.want)
		}
	}
}

func TestRetryDelayBounds(t *testing.T) {
	base := 100 * time.Millisecond
	max := 800 * time.Millisecond
	for attempt := 0; attempt < 8; attempt++ {
		for i := 0; i < 50; i++ {
			d := retryDelay(attempt, base, max)
			if d < base/2 {
				t.Errorf("attempt %d: delay %v below lower bound", attempt, d)
			}
			if d > max {
				t.Errorf("attempt %d: delay %v above cap %v", attempt, d, max)
			}
		}
	}
	if d := retryDelay(10, base, max); d > max {
		t.Errorf("exponential growth must be capped at %v, got %v", max, d)
	}
}
