package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/karldreher/gitstats/server/internal/handlers"
	"github.com/karldreher/gitstats/server/internal/metrics"
	"github.com/karldreher/gitstats/server/internal/persistence"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	store, err := persistence.FromEnv()
	if err != nil {
		log.Fatalf("persistence config error: %v", err)
	}

	if store != nil {
		saved, err := store.Load()
		if err != nil {
			log.Printf("warn: could not restore persisted metrics: %v", err)
		} else {
			for key, count := range saved {
				parts := strings.SplitN(key, "|", 4)
				if len(parts) == 4 {
					metrics.CommitsTotal.WithLabelValues(parts[0], parts[1], parts[2], parts[3]).Add(count)
				}
			}
		}
	}

	http.Handle("/", http.HandlerFunc(handlers.GetRoot))
	http.Handle("/api/v1/commit", handlers.NewPostCommit(store))
	http.Handle("/healthz", http.HandlerFunc(handlers.Healthz))
	http.Handle("/readyz", http.HandlerFunc(handlers.Readyz))
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8000", nil)
}
