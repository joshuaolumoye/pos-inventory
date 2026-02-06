package handler

import (
	"net/http"
)

// HealthCheckHandler is a simple handler to verify the server is running and routing works
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
