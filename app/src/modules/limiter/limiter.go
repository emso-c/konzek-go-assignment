// Package limiter provides a *token bucket algorithm* based rate limiter.
// It limits the rate of incoming requests from clients by keeping track of the
// requests made during a certain time window, and allowing a maximum number of
// requests per that window.
//
// Usage:
// Before using the limiter, set the rate limit and the window size by setting the
// environment variables `HTTP_RATE_LIMIT` and `HTTP_RATE_LIMIT_WINDOW` respectively.
// Then initialize the limiter by calling the `Initialize` function once at the start
// of your application. Finally, use the `ExceedsLimit` function to check if a client
// has exceeded the rate limit, and `Increment` function to record a request made by
// a client.
//
// Example:
// Set the environment variables:
// export HTTP_RATE_LIMIT=10
// export HTTP_RATE_LIMIT_WINDOW=60
//
// Initialize the limiter:
// limiter.GetLimiter().Initialize()
//
// Check if a client has exceeded the rate limit:
//
//	if limiter.GetLimiter().ExceedsLimit(clientID) {
//	    http.Error(w, "Too many requests", http.StatusTooManyRequests)
//	    return
//	}
//
// Increment the request count for a client:
// limiter.GetLimiter().Increment(clientID)
package limiter

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Limiter represents a token bucket algorithm based rate limiter.
type Limiter struct {
	UsageMap map[string]int
}

// NewLimiter returns a new Limiter instance.
func NewLimiter() *Limiter {
	return &Limiter{
		UsageMap: make(map[string]int),
	}
}

// ExceedsLimit checks if the number of requests made by a client exceeds the rate limit.
// It returns true if the limit is exceeded, false otherwise.
// The identifier parameter represents the client identifier (e.g. IP address, user ID, etc.).
func (l *Limiter) ExceedsLimit(identifier string) bool {
	limitStr := os.Getenv("HTTP_RATE_LIMIT")
	if limitStr == "" {
		log.Panic("HTTP_RATE_LIMIT is not set")
	}
	limit, _ := strconv.Atoi(limitStr)
	return l.UsageMap[identifier] > limit
}

// Increment increments the request count for a client.
// The identifier parameter represents the client identifier (e.g. IP address, user ID, etc.).
// Make sure the identifier is unique for each client.
func (l *Limiter) Increment(identifier string) {
	l.UsageMap[identifier]++
}

// Initialize initializes the limiter by setting a periodic task to reset
// the usage map after a certain time window.
// The time window and rate limit are read from the environment variables
// `HTTP_RATE_LIMIT_WINDOW` and `HTTP_RATE_LIMIT` respectively.
func (l *Limiter) Initialize() {
	windowStr := os.Getenv("HTTP_RATE_LIMIT_WINDOW")
	if windowStr == "" {
		log.Fatal("HTTP_RATE_LIMIT_WINDOW is not set")
	}
	window, err := strconv.Atoi(windowStr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(time.Duration(window) * time.Second)
			l.UsageMap = make(map[string]int)
		}
	}()
}

// l is the shared singleton instance of the Limiter.
var l *Limiter

// GetLimiter returns the shared singleton instance of the Limiter.
func GetLimiter() *Limiter {
	if l == nil {
		l = NewLimiter()
	}
	return l
}
