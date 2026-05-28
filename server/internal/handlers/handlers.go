package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/karldreher/gitstats/server/internal/metrics"
	"github.com/karldreher/gitstats/server/internal/persistence"
)

var conventionalRE = regexp.MustCompile(
	`^(feat|fix|docs|style|refactor|perf|test|chore)(\(.+\))?: .{1,50}`,
)
var commitTypeRE = regexp.MustCompile(
	`^(feat|fix|docs|style|refactor|perf|test|chore)`,
)

type CommitRequest struct {
	Commit string `json:"commit"`
	Repo   string `json:"repo"`
	Author string `json:"author"`
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// NewPostCommit returns a handler with an optional persistence store injected.
func NewPostCommit(store persistence.StateStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != os.Getenv("API_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req CommitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("decode error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Commit == "" || req.Repo == "" || req.Author == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		conventional := "false"
		commitType := "unknown"
		if conventionalRE.MatchString(req.Commit) {
			conventional = "true"
			if m := commitTypeRE.FindString(req.Commit); m != "" {
				commitType = m
			}
		}

		metrics.CommitsTotal.WithLabelValues(req.Repo, req.Author, commitType, conventional).Inc()

		if store != nil {
			key := persistence.LabelKey(req.Repo, req.Author, commitType, conventional)
			if err := store.Increment(key); err != nil {
				log.Println("persistence error:", err)
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}
