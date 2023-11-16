package store

import (
	"database/sql"
	"log"

	uuid "github.com/satori/go.uuid"
	"sigs.k8s.io/ggexample/models"
)

type QuizStore interface {
	Create() (*models.CreateQuizResponse, error)
	Exists(id string) (bool, error)
	DeleteByID(id string) error
}

type quizStore struct {
	db *sql.DB
}

func NewQuizStore(db *sql.DB) *quizStore {
	return &quizStore{
		db: db,
	}
}

func (s *quizStore) InitQuizRelation() error {
	return s.createQuizesTable()
}

func (s *quizStore) createQuizesTable() error {
	query := `create table if not exists quiz_sessions (
		id text primary key
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *quizStore) Create() (*models.CreateQuizResponse, error) {
	quizId := uuid.NewV4().String()

	query := `insert into quiz_sessions
	(id)
	values ($1)`

	_, err := s.db.Query(
		query,
		quizId,
	)

	if err != nil {
		return nil, err
	}

	return &models.CreateQuizResponse{
		ID: quizId,
	}, nil
}

func (s *quizStore) Exists(sessionID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("select EXISTS(select 1 from quiz_sessions where id = $1)", sessionID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *quizStore) DeleteByID(id string) error {
	_, err := s.db.Query("delete from quiz_sessions where id = $1", id)
	if err != nil {
		log.Println(err)
	}
	return err
}
