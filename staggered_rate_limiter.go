package main

import (
	"sync"
	"time"
)

type StaggeredRateLimiter struct {
	maxRequests int           // Maximum number of requests allowed
	timeWindow  time.Duration // Time window for the rate limit
	requests    []time.Time   // Timestamps of recent requests
	mu          sync.Mutex    // Mutex for thread safety
}

// NewStaggeredRateLimiter initializes a new StaggeredRateLimiter.
// NewStaggeredRateLimiter creates a new instance of StaggeredRateLimiter.
// It limits the number of requests to a specified maximum within a given time window.
//
// Parameters:
//   - maxRequests: The maximum number of requests allowed within the time window.
//   - timeWindow: The duration of the time window during which the requests are counted.
//
// Returns:
//
//	A pointer to a StaggeredRateLimiter instance configured with the specified limits.
func NewStaggeredRateLimiter(maxRequests int, timeWindow time.Duration) *StaggeredRateLimiter {
	return &StaggeredRateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// AllowRequest checks if a request is allowed under the staggered rate limit.
func (rl *StaggeredRateLimiter) AllowRequest() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Remove timestamps that are outside the time window
	cutoff := now.Add(-rl.timeWindow)
	for len(rl.requests) > 0 && rl.requests[0].Before(cutoff) {
		rl.requests = rl.requests[1:]
	}

	// Check if a new request can be allowed
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}
