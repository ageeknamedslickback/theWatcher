package v1

import (
	"fmt"
	"net/http"
)

var serverIsHealthy = func() bool { return true }

// HealthHandler checks if the server is healthy or not
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store")
		if serverIsHealthy() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"status":"ok"}`)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"status":"error"}`)
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "method is not allowed")
		return
	}
}
