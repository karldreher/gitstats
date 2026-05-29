package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/karldreher/gitstats/internal/ghclient"
	"github.com/karldreher/gitstats/internal/metrics"
	"github.com/karldreher/gitstats/internal/persistence"
	"github.com/karldreher/gitstats/internal/poller"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	client, err := ghclient.FromEnv()
	if err != nil {
		log.Fatalf("github config error: %v", err)
	}

	store, err := persistence.FromEnv()
	if err != nil {
		log.Fatalf("persistence config error: %v", err)
	}
	if store != nil {
		restoreMetrics(store)
	}

	go poller.Run(context.Background(), client, store)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if poller.IsReady() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8000", nil)
}

func restoreMetrics(store persistence.StateStore) {
	saved, err := store.Load()
	if err != nil {
		log.Printf("warn: could not restore persisted metrics: %v", err)
		return
	}
	for key, count := range saved {
		if strings.HasPrefix(key, "__") {
			continue // skip internal keys (e.g. __last_polled_at)
		}
		parts := strings.SplitN(key, "|", 4)
		if len(parts) == 4 {
			metrics.CommitsTotal.WithLabelValues(parts[0], parts[1], parts[2], parts[3]).Add(count)
		}
	}
}
