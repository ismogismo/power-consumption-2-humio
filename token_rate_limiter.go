package main

import (
	"sync"
	"time"
)

type TokenRateLimiter struct {
	maxRequests int           // Maximum number of requests allowed
	timeWindow  time.Duration // Time window for the rate limit
	tokens      float64       // Current number of tokens
	lastRefill  time.Time     // Last time tokens were refilled
	mu          sync.Mutex    // Mutex for thread safety
}

// NewTokenRateLimiter initializes a new RateLimiter.
func NewTokenRateLimiter(maxRequests int, timeWindow time.Duration) *TokenRateLimiter {
	return &TokenRateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		tokens:      float64(maxRequests),
		lastRefill:  time.Now(),
	}
}

// refillTokens refills tokens based on the elapsed time since the last refill.
func (rl *TokenRateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	newTokens := elapsed * (float64(rl.maxRequests) / rl.timeWindow.Seconds())
	if rl.tokens < float64(rl.maxRequests) {
		rl.tokens = min(float64(rl.maxRequests), rl.tokens+newTokens)
	}
	rl.lastRefill = now
}

// AllowRequest checks if a request is allowed under the rate limit.
func (rl *TokenRateLimiter) AllowRequest() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refillTokens()
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// min returns the smaller of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
