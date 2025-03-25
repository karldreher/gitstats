package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type CommitRecord struct {
	Commit string
	Repo   string
}

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatalf("Unable to ping database: %v\n", pingErr)
	}
	log.Println("Connected to database")

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
