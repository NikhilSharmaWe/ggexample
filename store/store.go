package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Dependency struct {
	QuestionStore QuestionStore
	QuizStore     QuizStore
}

func NewDB() (*sql.DB, error) {
	connStr := "user=miyamoto dbname=quiz password=1234 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
