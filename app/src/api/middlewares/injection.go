// Package middlewares provides protection against SQL injection attacks.
//
// Usage:
// Use the SQLInjectionMiddleware function as a middleware in your HTTP handlers to protect against SQL injection attacks.
//
// Example:
//
// http.Handle("/api/tasks", middlewares.SQLInjectionMiddleware()(http.HandlerFunc(handler)))
//
// The middleware checks for potential SQL injection patterns in the URL parameters, request body, and form data.
// If a potential SQL injection pattern is detected, the middleware returns a 400 Bad Request error.
package middlewares

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/emso-c/konzek-go-assignment/src/api/responses"
	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
)

// isSQLInjection checks if the provided string contains potential SQL injection patterns.
func isSQLInjection(value string) bool {
	patterns := []string{
		"SELECT",
		"DELETE",
		"DROP TABLE",
		"INSERT INTO",
		"UPDATE",
		"1=1",
		"OR 1=1",
		"AND 1=1",
		"OR '1'='1",
		"AND '1'='1",
		"OR 1=1--",
		"AND 1=1--",
		"OR '1'='1--",
		"AND '1'='1--",
	}
	for _, pattern := range patterns {
		if strings.Contains(value, pattern) {
			return true
		}
	}
	return false
}

// SQLInjectionMiddleware returns a middleware that protects against SQL injection attacks.
func SQLInjectionMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get logger
			var logger = logger.GetLogger()

			// Get request parameters
			params := r.URL.Query()
			for _, values := range params {
				for _, value := range values {
					if isSQLInjection(value) {
						responses.Error(w, http.StatusBadRequest, "Potential SQL Injection Detected")
						logger.Error("Potential SQL Injection Detected in URL parameter: " + value)
						return
					}
				}
			}

			// Get request body
			if r.Body != nil {
				maxBodySize := int64(1 << 20) // 1 MB
				body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
				if err != nil {
					responses.Error(w, http.StatusInternalServerError, "Internal Server Error")
					logger.Error("Error reading request body: " + err.Error())
					return
				}
				defer r.Body.Close()

				// Check for SQL injection in the request body
				if isSQLInjection(string(body)) {
					responses.Error(w, http.StatusBadRequest, "Potential SQL Injection Detected")
					logger.Error("Potential SQL Injection Detected in request body")
					return
				}

				// Reset the body to its original state
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

			// Get request form data
			if err := r.ParseForm(); err != nil {
				// Ignore error if there is no form data
				if !errors.Is(err, http.ErrNotMultipart) {
					responses.Error(w, http.StatusInternalServerError, "Internal Server Error")
					logger.Error("Error parsing form data: " + err.Error())
					return
				}
			}

			// Check for SQL injection in form data
			if len(r.PostForm) > 0 {
				for _, values := range r.PostForm {
					for _, value := range values {
						if isSQLInjection(value) {
							responses.Error(w, http.StatusBadRequest, "Potential SQL Injection Detected")
							logger.Error("Potential SQL Injection Detected in form data: " + value)
							return
						}
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
