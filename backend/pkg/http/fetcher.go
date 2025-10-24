// Package http provides HTTP utilities for fetching external data sources.
// Per AGENT_TECH_SPEC.md §9: retry with backoff, circuit-break on repeated failures.
package http

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// FetcherConfig holds configuration for HTTP fetching.
type FetcherConfig struct {
	MaxRetries     int           // Maximum number of retries (default: 3)
	InitialBackoff time.Duration // Initial backoff duration (default: 500ms)
	MaxBackoff     time.Duration // Maximum backoff duration (default: 30s)
	Timeout        time.Duration // HTTP request timeout (default: 30s)
	UserAgent      string        // User-Agent header
}

// DefaultConfig returns sensible defaults for fetching public data.
func DefaultConfig() FetcherConfig {
	return FetcherConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Timeout:        30 * time.Second,
		UserAgent:      "Corporate-Energy-Platform/1.0 (Educational Project)",
	}
}

// Fetcher performs HTTP requests with retry logic.
type Fetcher struct {
	config FetcherConfig
	client *http.Client
}

// NewFetcher creates a new HTTP fetcher with the given configuration.
func NewFetcher(config FetcherConfig) *Fetcher {
	return &Fetcher{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Fetch performs an HTTP GET request with exponential backoff retry.
// Returns the response body as an io.ReadCloser (caller must close).
func (f *Fetcher) Fetch(url string) (io.ReadCloser, error) {
	var lastErr error

	for attempt := 0; attempt <= f.config.MaxRetries; attempt++ {
		// Apply backoff before retry (skip on first attempt)
		if attempt > 0 {
			backoff := f.calculateBackoff(attempt)
			time.Sleep(backoff)
		}

		// Create request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set User-Agent
		if f.config.UserAgent != "" {
			req.Header.Set("User-Agent", f.config.UserAgent)
		}

		// Execute request
		resp, err := f.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed (attempt %d/%d): %w",
				attempt+1, f.config.MaxRetries+1, err)
			continue
		}

		// Check status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s (attempt %d/%d)",
				resp.StatusCode, url, attempt+1, f.config.MaxRetries+1)
			continue
		}

		// Success
		return resp.Body, nil
	}

	// All retries exhausted
	return nil, fmt.Errorf("fetch failed after %d attempts: %w",
		f.config.MaxRetries+1, lastErr)
}

// calculateBackoff computes exponential backoff duration.
// Formula: min(InitialBackoff * 2^attempt, MaxBackoff)
func (f *Fetcher) calculateBackoff(attempt int) time.Duration {
	backoff := float64(f.config.InitialBackoff) * math.Pow(2, float64(attempt-1))
	if backoff > float64(f.config.MaxBackoff) {
		backoff = float64(f.config.MaxBackoff)
	}
	return time.Duration(backoff)
}
