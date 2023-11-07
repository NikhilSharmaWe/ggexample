package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"sigs.k8s.io/ggexample/models"
)

type QuestionStore interface {
	Create(models.CreateQuestionRequest) error
	GetByID(int) (*models.GetQuestionResponse, error)
	DeleteByID(id int) error
}

type questionStore struct {
	db *sql.DB
}

func NewQuestionStore() (*questionStore, error) {
	connStr := "user=miyamoto dbname=quiz password=1234 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &questionStore{
		db: db,
	}, nil
}

func (s *questionStore) InitQuestionRelation() error {
	return s.createQuestionsTable()
}

func (s *questionStore) createQuestionsTable() error {
	query := `create table if not exists questions (
		id serial primary key,
		text varchar(500),
		options text[],
		answer text
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *questionStore) Create(q models.CreateQuestionRequest) error {
	query := `insert into questions 
	(text, options, answer)
	values ($1, $2, $3)`

	optionsArray := pq.Array(q.Options)

	_, err := s.db.Query(
		query,
		q.Text,
		optionsArray,
		q.Answer,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *questionStore) DeleteByID(id int) error {
	_, err := s.db.Query("delete from questions where id = $1", id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *questionStore) GetByID(id int) (*models.GetQuestionResponse, error) {
	rows, err := s.db.Query("select text, options from questions where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuestion(rows)
	}

	return nil, fmt.Errorf("question %d not found", id)
}

func scanIntoQuestion(rows *sql.Rows) (*models.GetQuestionResponse, error) {
	q := new(models.GetQuestionResponse)

	var pqArray pq.StringArray

	err := rows.Scan(
		&q.Text,
		&pqArray,
	)

	q.Options = pqArray
	return q, err
}
