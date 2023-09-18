package utils

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler returns the status of the server
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server is up and running!")
}
