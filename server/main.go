package main

import (
	"log"
	"net/http"

	gitstats "github.com/karldreher/gitstats/server/src"
)

func main() {
	http.HandleFunc("/", gitstats.GetRoot)
	http.HandleFunc("/api/v1/commit", gitstats.PostCommit)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
