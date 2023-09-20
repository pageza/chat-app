// Package utils contains utility functions that can be used across the application.

package utils

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse is a utility function to send a JSON response.
// It takes an http.ResponseWriter, a status code, and a payload (interface{} type).
// The function sets the Content-Type to "application/json" and writes the JSON-encoded payload to the response.
func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	// Set the Content-Type to "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Set the status code for the response
	w.WriteHeader(statusCode)

	// Encode the payload to JSON and write it to the response
	json.NewEncoder(w).Encode(payload)
}
