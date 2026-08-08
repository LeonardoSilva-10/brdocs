// Package httpclient provides shared HTTP utilities for the brdocs API clients.
package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultTimeout is the default HTTP client timeout.
const DefaultTimeout = 30 * time.Second

// NewClient creates a new HTTP client with sensible defaults.
func NewClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &http.Client{
		Timeout: timeout,
	}
}

// ReadBody safely reads and closes an HTTP response body.
func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// CheckStatus verifies that the response status code is in the expected range.
func CheckStatus(resp *http.Response, expected ...int) error {
	for _, code := range expected {
		if resp.StatusCode == code {
			return nil
		}
	}
	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// DoWithContext performs an HTTP request with the given context.
func DoWithContext(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}
