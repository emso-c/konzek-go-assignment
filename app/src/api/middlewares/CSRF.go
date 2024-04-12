// Package middlewares provides protection against Cross-Site Request Forgery (CSRF) attacks.
//
// Usage:
// Use the CSRFMiddleware function as a middleware in your HTTP handlers to protect against CSRF attacks.
//
// Example:
//
// http.Handle("/api/tasks", middlewares.CSRFMiddleware()(http.HandlerFunc(handler)))
//
// The middleware generates a CSRF token for the current request, and sets it in the response header and cookie.
// It checks the CSRF token in the request header or cookie for POST, PUT, and DELETE requests.
//
// If the CSRF token is missing or invalid, the middleware returns a 403 Forbidden error.
package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const (
	csrfTokenLength = 32
	csrfHeaderName  = "X-CSRF-Token"
	csrfCookieName  = "csrf_token"
	csrfMaxAge      = 3600 // 1 hour
)

// generateCSRFToken generates a CSRF token for the current request.
func generateCSRFToken(r *http.Request) (string, error) {
	// Attempt to read from cookie
	cookie, err := r.Cookie(csrfCookieName)
	if err == nil && cookie != nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	// Generate new token
	token := make([]byte, csrfTokenLength)
	_, err = rand.Read(token)
	if err != nil {
		return "", err
	}

	// Encode token
	return base64.URLEncoding.EncodeToString(token), nil
}

// CSRFMiddleware returns a middleware that protects against Cross-Site Request Forgery (CSRF) attacks.
func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate CSRF token if not present
			csrfToken, err := generateCSRFToken(r)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set CSRF token in response header
			w.Header().Set(csrfHeaderName, csrfToken)

			// Set CSRF token in cookie
			http.SetCookie(w, &http.Cookie{
				Name:     csrfCookieName,
				Value:    csrfToken,
				HttpOnly: true,
				MaxAge:   csrfMaxAge,
				SameSite: http.SameSiteStrictMode,
			})

			// Check CSRF token in request header or cookie
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
				clientToken := r.Header.Get(csrfHeaderName)
				if clientToken == "" {
					clientToken = r.FormValue(csrfCookieName)
				}
				if clientToken != csrfToken {
					http.Error(w, "CSRF Token Invalid", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
