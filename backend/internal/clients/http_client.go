package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// BaseHTTPClient implements HTTPClient with retry logic and rate limiting
type BaseHTTPClient struct {
	client     *http.Client
	limiter    *rate.Limiter
	maxRetries int
	retryDelay time.Duration
}

// NewBaseHTTPClient creates a new HTTP client with retry logic
func NewBaseHTTPClient(config ClientConfig) *BaseHTTPClient {
	// Create rate limiter
	limiter := rate.NewLimiter(
		rate.Limit(config.RateLimit.RequestsPerSecond),
		config.RateLimit.BurstSize,
	)

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &BaseHTTPClient{
		client:     httpClient,
		limiter:    limiter,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// Do executes an HTTP request with retry logic and rate limiting
func (c *BaseHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Clone request for retry
		reqClone := c.cloneRequest(req)
		
		resp, err := c.client.Do(reqClone)
		if err != nil {
			lastErr = err
			if attempt < c.maxRetries {
				c.waitForRetry(req.Context(), attempt)
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, err)
		}

		// Check if we should retry based on status code
		if c.shouldRetry(resp.StatusCode) && attempt < c.maxRetries {
			resp.Body.Close()
			lastErr = fmt.Errorf("received status code %d", resp.StatusCode)
			c.waitForRetry(req.Context(), attempt)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

// Get performs a GET request
func (c *BaseHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post performs a POST request
func (c *BaseHTTPClient) Post(url, contentType string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// cloneRequest creates a copy of the request for retry
func (c *BaseHTTPClient) cloneRequest(req *http.Request) *http.Request {
	clone := req.Clone(req.Context())
	
	// Clone body if present
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		clone.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	return clone
}

// shouldRetry determines if a request should be retried based on status code
func (c *BaseHTTPClient) shouldRetry(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		 http.StatusInternalServerError,
		 http.StatusBadGateway,
		 http.StatusServiceUnavailable,
		 http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// waitForRetry implements exponential backoff with jitter
func (c *BaseHTTPClient) waitForRetry(ctx context.Context, attempt int) {
	// Exponential backoff: base_delay * 2^attempt
	delay := c.retryDelay * time.Duration(1<<uint(attempt))
	
	// Cap the delay at 30 seconds
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}

// ParseResponse parses HTTP response body into a struct
func ParseResponse[T any](resp *http.Response, target *T) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("API error [%s]: %s", errResp.Code, errResp.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}