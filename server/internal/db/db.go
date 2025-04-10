package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type CommitRecord struct {
	Commit string
	Repo   string
}

var DB *sql.DB

func Connect() {
	var err error
	for i := range 5 {
		DB, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
		if err == nil {
			pingErr := DB.Ping()
			if pingErr == nil {
				log.Println("Connected to database")
				return
			}
			log.Printf("Ping attempt %d failed: %v\n", i+1, pingErr)
			// Wait i * 2(i^2) seconds before retrying
			time.Sleep(time.Duration(2*i*i) * time.Second)
		} else {
			log.Printf("Connection attempt %d failed: %v\n", i+1, err)
		}
	}
	log.Fatalf("Unable to connect to database after 5 attempts: %v\n", err)

}

func InsertCommit(commit CommitRecord) error {
	_, err := DB.Exec("INSERT INTO commits (commit, repo) VALUES ($1, $2)", commit.Commit, commit.Repo)
	if err != nil {
		return err
	}
	return nil
}

func DBReady() bool {
	err := DB.Ping()
	return err == nil
}
