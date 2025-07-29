package external

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens   int
	maxTokens int
	interval time.Duration
	mu       sync.Mutex
	ticker   *time.Ticker
}

func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:    maxRequests,
		maxTokens: maxRequests,
		interval:  interval,
	}

	// Refill tokens periodically
	rl.ticker = time.NewTicker(interval / time.Duration(maxRequests))
	go func() {
		for range rl.ticker.C {
			rl.mu.Lock()
			if rl.tokens < rl.maxTokens {
				rl.tokens++
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			rl.mu.Lock()
			if rl.tokens > 0 {
				rl.tokens--
				rl.mu.Unlock()
				return nil
			}
			rl.mu.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// Stop cleans up the rate limiter
func (rl *RateLimiter) Stop() {
	if rl.ticker != nil {
		rl.ticker.Stop()
	}
}