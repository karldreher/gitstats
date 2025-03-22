package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/karldreher/gitstats/server/internal/db"
)

type Commit struct {
	Commit string `json:"commit"`
	Repo   string `json:"repo"`
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func PostCommit(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("x-api-key")
	// TODO, a more abstract API key (in the db?)
	if h != os.Getenv("API_KEY") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	j := json.NewDecoder(r.Body)
	var commit Commit
	err := j.Decode(&commit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// A little redundant but fine for now, more strict on both sides is not bad
	c := db.CommitRecord{
		Commit: commit.Commit,
		Repo:   commit.Repo,
	}
	err = db.InsertCommit(c)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
