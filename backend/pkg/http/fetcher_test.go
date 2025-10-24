package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetcher_Fetch_Success(t *testing.T) {
	// Create test server that returns success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check User-Agent
		if ua := r.Header.Get("User-Agent"); ua == "" {
			t.Error("User-Agent header not set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	config := DefaultConfig()
	fetcher := NewFetcher(config)

	body, err := fetcher.Fetch(server.URL)
	if err != nil {
		t.Fatalf("Fetch() error = %v, want nil", err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if string(data) != "test data" {
		t.Errorf("Body = %q, want %q", string(data), "test data")
	}
}

func TestFetcher_Fetch_Retry(t *testing.T) {
	attempts := 0

	// Server that fails twice, then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success after retry"))
	}))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 3
	config.InitialBackoff = 10 * time.Millisecond // Fast for tests

	fetcher := NewFetcher(config)

	body, err := fetcher.Fetch(server.URL)
	if err != nil {
		t.Fatalf("Fetch() error = %v, want nil", err)
	}
	defer body.Close()

	if attempts != 3 {
		t.Errorf("Attempts = %d, want 3", attempts)
	}
}

func TestFetcher_Fetch_Exhausted(t *testing.T) {
	// Server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 2
	config.InitialBackoff = 10 * time.Millisecond

	fetcher := NewFetcher(config)

	_, err := fetcher.Fetch(server.URL)
	if err == nil {
		t.Error("Fetch() error = nil, want error after exhausted retries")
	}
}

func TestFetcher_CalculateBackoff(t *testing.T) {
	config := DefaultConfig()
	config.InitialBackoff = 100 * time.Millisecond
	config.MaxBackoff = 1 * time.Second

	fetcher := NewFetcher(config)

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 100 * time.Millisecond},  // 100 * 2^0
		{2, 200 * time.Millisecond},  // 100 * 2^1
		{3, 400 * time.Millisecond},  // 100 * 2^2
		{4, 800 * time.Millisecond},  // 100 * 2^3
		{5, 1 * time.Second},         // capped at MaxBackoff
		{10, 1 * time.Second},        // capped at MaxBackoff
	}

	for _, tt := range tests {
		got := fetcher.calculateBackoff(tt.attempt)
		if got != tt.want {
			t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}
