package gitstats

import (
	"encoding/json"
	"log"
	"net/http"
)

type Commit struct {
	Commit string `json:"commit"`
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func PostCommit(w http.ResponseWriter, r *http.Request) {

	j := json.NewDecoder(r.Body)
	var commit Commit
	err := j.Decode(&commit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
