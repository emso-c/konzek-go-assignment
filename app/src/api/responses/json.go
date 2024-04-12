// This package contains the response functions for the API
//
// # This file contains the JSON response function
//
// Usage:
// Use the JSON function to return a JSON response
//
// Example:
//
// err := responses.JSON(w, http.StatusOK, tasks)
//
//	if err != nil {
//		panic(err)
//	}
package responses

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	return nil
}
