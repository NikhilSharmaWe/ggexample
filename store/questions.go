package store

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type QuestionStore interface {
	Create(*Question) error
	DeleteByID(id int) error
	GetByID(int) (*Question, error)
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

func (s *questionStore) Create(q *Question) error {
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
	return err
}

func (s *questionStore) GetByID(id int) (*Question, error) {
	rows, err := s.db.Query("select * from questions where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuestion(rows)
	}

	return nil, fmt.Errorf("question %d not found", id)
}

func scanIntoQuestion(rows *sql.Rows) (*Question, error) {
	q := new(Question)

	var pqArray pq.StringArray

	err := rows.Scan(
		&q.ID,
		&q.Text,
		&pqArray,
		&q.Answer,
	)

	q.Options = pqArray
	return q, err
}
