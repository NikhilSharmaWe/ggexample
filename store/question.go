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
	Create(models.CreateQuestionRequest) error
	GetByID(int) (*models.GetQuestionResponse, error)
	GetRandomQuestions(int) (*models.GetQuestionsResponse, error)
	GetRandomQuestionIds(int) ([]int, error)
	CheckAnswer(id int, answer string) (bool, error)
	DeleteByID(id int) error
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
		fmt.Println("--------------------")
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

		question, err := s.GetByID(randomID)
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

func (s *questionStore) GetRandomQuestionIds(count int) ([]int, error) {
	totalQuestions, err := s.getTotalQuestionCount()
	if err != nil {
		return nil, err
	}

	if totalQuestions < count {
		return nil, fmt.Errorf("not enough questions available in the database")
	}

	var questionsIds []int
	selectedQuestions := make(map[int]struct{})

	for len(questionsIds) < count {
		randomID := rand.Intn(totalQuestions) + 1

		if _, ok := selectedQuestions[randomID]; ok {
			continue
		}

		_, err := s.GetByID(randomID)
		if err != nil {
			continue
		}

		questionsIds = append(questionsIds, randomID)
		selectedQuestions[randomID] = struct{}{}
	}

	return questionsIds, nil
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

func (s *questionStore) getTotalQuestionCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT count(*) FROM questions").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
