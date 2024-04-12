// Package middlewares provides HTTP middlewares for the API.
//
// Usage:
// Use the CORSMiddleware function as a middleware in your HTTP handlers to enable CORS.
// The middleware retrieves the origin from the request header, and sets the Access-Control-Allow-Origin
// header in the response if the origin is allowed. It also sets the Access-Control-Allow-Methods header
// in the response for a Preflighted OPTIONS request.
//
// Example:
//
// http.Handle("/api/tasks", middlewares.CORSMiddleware()(http.HandlerFunc(handler)))
//
// The middleware returns a 400 Bad Request error if the remote address is missing.
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
)

func _contains(s string, e string) bool {
	for _, a := range strings.Split(s, ",") {
		if a == e {
			return true
		}
	}
	return false
}

// CORSMiddleware returns a middleware that enables CORS.
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.GetLogger()
			origin := r.Header.Get("Origin")
			allowedOrigins := strings.Split(",", os.Getenv("HTTP_ALLOWED_ORIGINS"))
			if origin != "" {
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			// Check if the request method is allowed
			if os.Getenv("HTTP_ALLOWED_METHODS") != "*" {
				if !_contains(r.Method, os.Getenv("HTTP_ALLOWED_METHODS")) {
					http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
					logger.Error(fmt.Sprintf("Method Not Allowed: %s", r.Method))
					return
				}
			}

			// Stop here for a Preflighted OPTIONS request
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", os.Getenv("HTTP_ALLOWED_METHODS"))
				w.Header().Set("Access-Control-Allow-Headers", os.Getenv("HTTP_ALLOWED_HEADERS"))
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
