package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"sigs.k8s.io/ggexample/models"
)

type QuizStore interface {
	Create(models.CreateQuizRequest) error
	GetByID(string) (*models.GetQuizResponse, error)
	// GetRandomQuestionsID(int) (*models.GetQuestionsResponse, error)
	// CheckAnswer(id int, answer string) (bool, error)
	DeleteByID(id string) error
	UpdateQuiz(models.UpdateQuizRequest) error
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
	query := `create table if not exists quizes (
		id text primary key,
		questionids INT[],
		progress text[]
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *quizStore) Create(q models.CreateQuizRequest) error {
	query := `insert into quizes
	(id, questionids, progress)
	values ($1, $2, $3)`

	questionids := pq.Array(q.QuestionIDs)
	progressArray := pq.Array([]string{})

	_, err := s.db.Query(
		query,
		q.ID,
		questionids,
		progressArray,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *quizStore) DeleteByID(id string) error {
	_, err := s.db.Query("delete from quizes where id = $1", id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *quizStore) GetByID(id string) (*models.GetQuizResponse, error) {
	rows, err := s.db.Query("select id, questionids, progress from quizes where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoQuiz(rows)
	}

	return nil, fmt.Errorf("quiz %s not found", id)
}

func (s *quizStore) UpdateQuiz(req models.UpdateQuizRequest) error {
	query := `
        update quizes
        set progress[$2] = $3
        where id = $1
    `

	_, err := s.db.Exec(query, req.ID, req.Index, req.Answer)

	return err
}

// func (s *questionStore) GetRandomQuestions(count int) (*models.GetQuestionsResponse, error) {
// 	totalQuestions, err := s.getTotalQuestionCount()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if totalQuestions < count {
// 		return nil, fmt.Errorf("not enough questions available in the database")
// 	}

// 	selectedQuestions := make(map[int]struct{})
// 	questions := make([]*models.GetQuestionResponse, 0)
// 	for len(questions) < count {
// 		randomID := rand.Intn(totalQuestions) + 1

// 		if _, ok := selectedQuestions[randomID]; ok {
// 			continue
// 		}

// 		question, err := s.GetQuestionByID(randomID)
// 		if err != nil {
// 			continue
// 		}

// 		questions = append(questions, question)
// 		selectedQuestions[randomID] = struct{}{}
// 	}

// 	return &models.GetQuestionsResponse{
// 		Questions: questions,
// 	}, nil
// }

// func (s *questionStore) CheckAnswer(id int, answer string) (bool, error) {
// 	q, err := s.GetQuestionByID(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	return answer == q.Answer, nil
// }

func scanIntoQuiz(rows *sql.Rows) (*models.GetQuizResponse, error) {
	q := new(models.GetQuizResponse)

	var questionIds pq.Int64Array
	var pqArray pq.StringArray

	err := rows.Scan(
		&q.ID,
		&questionIds,
		&pqArray,
	)

	q.QuestionIDs = make([]int, len(questionIds))
	for i, id := range questionIds {
		q.QuestionIDs[i] = int(id)
	}
	// fmt.Println(questionIds)

	q.Progress = pqArray
	// q.QuestionIDs = /
	return q, err
}

// func (s *questionStore) getTotalQuestionCount() (int, error) {
// 	var count int
// 	err := s.db.QueryRow("SELECT count(*) FROM questions").Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }
