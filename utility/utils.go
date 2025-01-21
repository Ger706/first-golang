package utility

import (
	"net/http"
)

// SetJSONHeader sets the response header for JSON content type
func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
