package handlers

import (
	"net/http"

	"github.com/karldreher/gitstats/server/internal/poller"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz returns 200 once the first poll has completed, 503 until then.
func Readyz(w http.ResponseWriter, r *http.Request) {
	if poller.IsReady() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
