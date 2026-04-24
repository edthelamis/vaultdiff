package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newVersionedVaultServer returns a test server that serves version 1 on the
// first call and version 2 on subsequent calls.
func newVersionedVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	var calls int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := atomic.AddInt32(&calls, 1)
		if call == 1 {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"data":{"key":"v1"}}}`))
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"data":{"key":"v2"}}}`))
		}
	}))
}

func TestWatchSecret_DetectsChange(t *testing.T) {
	srv := newVersionedVaultServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := WatchOptions{Interval: 50 * time.Millisecond, MaxPolls: 2}
	ch := WatchSecret(ctx, client, "dev", "secret/data/app", 1, opts)

	var events []WatchEvent
	for e := range ch {
		events = append(events, e)
	}

	if len(events) == 0 {
		t.Fatal("expected at least one watch event")
	}
}

func TestWatchSecret_CancelStops(t *testing.T) {
	srv := newTestVaultServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	opts := WatchOptions{Interval: 20 * time.Millisecond}
	ch := WatchSecret(ctx, client, "dev", "secret/data/app", 1, opts)

	cancel()
	// Drain channel; it must close without blocking.
	for range ch {
	}
}

func TestWatchSecret_DefaultInterval(t *testing.T) {
	// Ensure a zero interval is replaced with the default (no panic).
	srv := newTestVaultServer(t)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "token")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	opts := WatchOptions{Interval: 0, MaxPolls: 1}
	ch := WatchSecret(ctx, client, "dev", "secret/data/app", 1, opts)
	for range ch {
	}
}
