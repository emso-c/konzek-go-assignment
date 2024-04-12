// Package middlewares provides HTTP middlewares for the API.
//
// Usage:
// Use the NotFoundMiddleware function as a middleware in your HTTP handlers to return a 404 Not Found
// error if the request path is not found.
//
// Example:
//
// http.Handle("/api", middlewares.NotFoundMiddleware())
//
// The middleware returns a 404 Not Found error if the request path is not found.
package middlewares

import (
	"encoding/json"
	"net/http"
)

// NotFoundMiddleware returns a 404 Not Found error if the request path is not found.
func NotFoundMiddleware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Page Not found"})
	})
}
