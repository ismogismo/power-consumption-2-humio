package main

import (
	"sync"
	"time"
)

type SlidingWindowRateLimiter struct {
	maxRequests int           // Maximum number of requests allowed
	timeWindow  time.Duration // Time window for the rate limit
	requests    map[int64]int // Map to store request counts per time window (key: timestamp in seconds)
	mu          sync.Mutex    // Mutex for thread safety
}

// NewSlidingWindowRateLimiter initializes a new SlidingWindowRateLimiter.
// NewSlidingWindowRateLimiter creates a new instance of SlidingWindowRateLimiter.
// It initializes the rate limiter with the specified maximum number of requests
// allowed (`maxRequests`) within a given time window (`timeWindow`).
//
// Parameters:
//   - maxRequests: The maximum number of requests allowed within the time window.
//   - timeWindow: The duration of the time window for rate limiting.
//
// Returns:
//
//	A pointer to a newly created SlidingWindowRateLimiter instance.
func NewSlidingWindowRateLimiter(maxRequests int, timeWindow time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		requests:    make(map[int64]int),
	}
}

// AllowRequest checks if a request is allowed under the sliding window rate limit.
func (rl *SlidingWindowRateLimiter) AllowRequest() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	currentWindow := now.Unix() // Current timestamp in seconds
	windowStart := now.Add(-rl.timeWindow).Unix()

	// Clean up old windows outside the time window
	for timestamp := range rl.requests {
		if timestamp < windowStart {
			delete(rl.requests, timestamp)
		}
	}

	// Calculate the total number of requests in the current sliding window
	totalRequests := 0
	for timestamp, count := range rl.requests {
		if timestamp >= windowStart {
			totalRequests += count
		}
	}

	// Check if the new request can be allowed
	if totalRequests < rl.maxRequests {
		rl.requests[currentWindow]++
		return true
	}

	return false
}
