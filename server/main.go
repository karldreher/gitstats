package main

import (
	"net/http"

	"github.com/karldreher/gitstats/server/internal/db"
	"github.com/karldreher/gitstats/server/internal/handlers"
)

func main() {
	db.Connect()
	http.Handle("/", http.HandlerFunc(handlers.GetRoot))
	http.Handle("/api/v1/commit", http.HandlerFunc(handlers.PostCommit))
	http.Handle("/healthz", http.HandlerFunc(handlers.Healthz))
	http.Handle("/readyz", http.HandlerFunc(handlers.Readyz))

	http.ListenAndServe(":8000", nil)
}
