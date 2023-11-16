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
	CheckAnswer(id int, answer string) (bool, error)
	DeleteByID(id int) error
	GetNextQuestion(sessionID string) (*models.GetQuestionResponse, error)
}

type questionStore struct {
	db *sql.DB
}

func NewQuestionStore(db *sql.DB) *questionStore {
	return &questionStore{
		db: db,
	}
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
	rows, err := s.db.Query("select id, text, options, answer from questions where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuestion(rows)
	}

	return nil, fmt.Errorf("question %d not found", id)
}

func (s *questionStore) CheckAnswer(id int, answer string) (bool, error) {
	q, err := s.GetByID(id)
	if err != nil {
		return false, err
	}

	return answer == q.Answer, nil
}

func scanIntoQuestion(rows *sql.Rows) (*models.GetQuestionResponse, error) {
	q := new(models.GetQuestionResponse)

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

func (s *questionStore) GetNextQuestion(sessionID string) (*models.GetQuestionResponse, error) {
	rows, err := s.db.Query(`
        select id, text, options, answer
        FROM questions
        WHERE id NOT IN (
            SELECT question_id
            FROM quiz_responses
            WHERE quiz_session_id = $1
        )
        ORDER BY RANDOM()
        LIMIT 1
    `, sessionID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuestion(rows)
	}

	return nil, fmt.Errorf("question not found")
}
