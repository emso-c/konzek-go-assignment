package responses

import (
	"encoding/json"
	"net/http"
)

func Plain(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}
