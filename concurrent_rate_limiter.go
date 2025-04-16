package main

import (
	"sync"
)

type ConcurrentRateLimiter struct {
	maxConcurrentRequests int        // Maximum number of concurrent requests allowed
	currentRequests       int        // Current number of active requests
	mu                    sync.Mutex // Mutex for thread safety
	cond                  *sync.Cond // Condition variable to manage waiting requests
}

// NewConcurrentRateLimiter initializes a new ConcurrentRateLimiter.
func NewConcurrentRateLimiter(maxConcurrentRequests int) *ConcurrentRateLimiter {
	limiter := &ConcurrentRateLimiter{
		maxConcurrentRequests: maxConcurrentRequests,
	}
	limiter.cond = sync.NewCond(&limiter.mu)
	return limiter
}

// Acquire blocks until a request can be processed.
func (rl *ConcurrentRateLimiter) Acquire() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Wait until there is room for a new request
	for rl.currentRequests >= rl.maxConcurrentRequests {
		rl.cond.Wait()
	}

	// Increment the count of active requests
	rl.currentRequests++
}

// Release signals that a request has finished processing.
func (rl *ConcurrentRateLimiter) Release() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Decrement the count of active requests
	rl.currentRequests--

	// Signal waiting goroutines that a slot is available
	rl.cond.Signal()
}
