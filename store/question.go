package store

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/lib/pq"
	"sigs.k8s.io/ggexample/models"
)

type QuestionStore interface {
	CreateQuestion(models.CreateQuestionRequest) error
	GetQuestionByID(int) (*models.GetQuestionResponse, error)
	GetRandomQuestions(int) (*models.GetQuestionsResponse, error)
	CheckAnswer(id int, answer string) (bool, error)
	DeleteQuestionByID(id int) error
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

func (s *questionStore) CreateQuestion(q models.CreateQuestionRequest) error {
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

func (s *questionStore) DeleteQuestionByID(id int) error {
	_, err := s.db.Query("delete from questions where id = $1", id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *questionStore) GetQuestionByID(id int) (*models.GetQuestionResponse, error) {
	rows, err := s.db.Query("select id, text, options, answer from questions where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuestion(rows)
	}

	return nil, fmt.Errorf("question %d not found", id)
}

func (s *questionStore) GetRandomQuestions(count int) (*models.GetQuestionsResponse, error) {
	totalQuestions, err := s.getTotalQuestionCount()
	if err != nil {
		return nil, err
	}

	if totalQuestions < count {
		return nil, fmt.Errorf("not enough questions available in the database")
	}

	selectedQuestions := make(map[int]struct{})
	questions := make([]*models.GetQuestionResponse, 0)
	for len(questions) < count {
		randomID := rand.Intn(totalQuestions) + 1

		if _, ok := selectedQuestions[randomID]; ok {
			continue
		}

		question, err := s.GetQuestionByID(randomID)
		if err != nil {
			continue
		}

		questions = append(questions, question)
		selectedQuestions[randomID] = struct{}{}
	}

	return &models.GetQuestionsResponse{
		Questions: questions,
	}, nil
}

func (s *questionStore) CheckAnswer(id int, answer string) (bool, error) {
	q, err := s.GetQuestionByID(id)
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

func (s *questionStore) getTotalQuestionCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT count(*) FROM questions").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
