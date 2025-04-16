package main

import (
	"time"
)

// Example usage
func main() {
	// tokenRateLimiter := NewTokenRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	// for i := 1; i <= 10; i++ {
	// 	if tokenRateLimiter.AllowRequest() {
	// 		println("Request", i, "allowed")
	// 	} else {
	// 		println("Request", i, "denied")
	// 	}
	// 	time.Sleep(500 * time.Millisecond)
	// }

	// sRateLimiter := NewStaggeredRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	// for i := 1; i <= 10; i++ {
	// 	if sRateLimiter.AllowRequest() {
	// 		println("Request", i, "allowed")
	// 	} else {
	// 		println("Request", i, "denied")
	// 	}
	// 	time.Sleep(200 * time.Millisecond) // Simulate staggered requests
	// }

	slidingRateLimiter := NewSlidingWindowRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	for i := 1; i <= 10; i++ {
		if slidingRateLimiter.AllowRequest() {
			println("Request", i, "allowed")
		} else {
			println("Request", i, "denied")
		}
		time.Sleep(400 * time.Millisecond) // Simulate staggered requests
	}
}
