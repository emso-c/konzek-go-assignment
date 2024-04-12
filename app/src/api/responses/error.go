// This package contains the response functions for the API
//
// # This file contains the Error response function
//
// Usage:
// Use the Error function to return an error response
//
// Example:
//
// err := responses.Error(w, http.StatusBadRequest, "Invalid request")
//
//	if err != nil {
//		panic(err)
//	}
package responses

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(map[string]map[string]interface{}{"error": {"code": code, "message": message}})
}
