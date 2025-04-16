package main

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rateLimiter := NewTokenRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	// Test that the first 5 requests are allowed
	for i := 1; i <= 5; i++ {
		if !rateLimiter.AllowRequest() {
			t.Errorf("Request %d was denied, but it should have been allowed", i)
		}
	}

	// Test that the 6th request is denied
	if rateLimiter.AllowRequest() {
		t.Error("Request 6 was allowed, but it should have been denied")
	}

	// Wait for the rate limiter to refill tokens
	time.Sleep(10 * time.Second)

	// Test that requests are allowed again after refill
	for i := 1; i <= 5; i++ {
		if !rateLimiter.AllowRequest() {
			t.Errorf("Request %d was denied after refill, but it should have been allowed", i)
		}
	}
}
