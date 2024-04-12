// Package middlewares provides HTTP middlewares for the API.
//
// Usage:
// Use the RateLimitMiddleware function as a middleware in your HTTP handlers to limit the
// number of requests from a single remote address. The middleware retrieves the remote address
// from the request header, or from the environment variable MOCK_REMOTE_ADDR if the environment
// variable ENV_LOCAL is set to "true". The middleware uses a limiter to count and limit the number
// of requests from each remote address.
//
// Example:
//
// Use the RateLimitMiddleware function as a middleware in your HTTP handlers:
//
// http.Handle("/api/tasks", middlewares.RateLimitMiddleware()(http.HandlerFunc(handler)))
//
// The middleware returns a 400 Bad Request error if the remote address is missing, and a 429 Too Many
// Requests error if the remote address has exceeded the rate limit.
//
// The GetRemoteAddr function retrieves the remote address from the request header, or from the
// environment variable MOCK_REMOTE_ADDR if the environment variable ENV_LOCAL is set to "true".
// If the remote address is missing, GetRemoteAddr returns an empty string.
//
// Example:
// Get the remote address from the request:
//
// remoteAddr := middlewares.GetRemoteAddr(r)
//
// You can use the remote address to identify the client making the request, or to implement custom
// rate limiting logic.
package middlewares

import (
	"net/http"
	"os"

	"github.com/emso-c/konzek-go-assignment/src/api/responses"
	"github.com/emso-c/konzek-go-assignment/src/modules/limiter"
)

// GetRemoteAddr retrieves the remote address from the request header, or from the environment
// if the environment variable ENV_LOCAL is set to "true".
func GetRemoteAddr(r *http.Request) string {
	remoteAddr := r.Header.Get("REMOTE_ADDR")
	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Forwarded-For")
	}
	// Use mock remote address for local testing
	if os.Getenv("ENV_LOCAL") == "true" {
		remoteAddr = os.Getenv("MOCK_REMOTE_ADDR")
	}
	return remoteAddr
}

// RateLimitMiddleware returns a middleware that limits the number of requests from a single remote address.
func RateLimitMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteAddr := GetRemoteAddr(r)
			if remoteAddr == "" {
				err := responses.Error(w, http.StatusBadRequest, "Bad request, missing remote address")
				if err != nil {
					panic(err)
				}
				return
			}

			l := limiter.GetLimiter()
			l.Increment(remoteAddr)
			if l.ExceedsLimit(remoteAddr) {
				err := responses.Error(w, http.StatusTooManyRequests, "Too many requests")
				if err != nil {
					panic(err)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
