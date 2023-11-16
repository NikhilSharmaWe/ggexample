package store

import (
	"database/sql"

	"sigs.k8s.io/ggexample/models"
)

type ResponseStore interface {
	Create(models.CreateResponseRequest) error
	IsQuizCompleted(sessionID string) (bool, error)
	GetQuizResult(sessionID string) (*models.QuizResultResponse, error)
}

type responseStore struct {
	db *sql.DB
}

func NewResponseStore(db *sql.DB) *responseStore {
	return &responseStore{
		db: db,
	}
}

func (s *responseStore) InitResponseRelation() error {
	return s.createResponseTable()
}

func (s *responseStore) createResponseTable() error {
	query := `create table if not exists quiz_responses (
		id SERIAL PRIMARY KEY,
    	quiz_session_id TEXT REFERENCES quiz_sessions(id),
    	question_id INT REFERENCES questions(id),
    	answer TEXT,
    	is_correct BOOLEAN
	)`

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	_, err := s.db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS unique_response ON quiz_responses(quiz_session_id, question_id)
	`)

	return err
}

func (s *responseStore) Create(q models.CreateResponseRequest) error {
	query := `insert into quiz_responses
	(quiz_session_id, question_id, answer, is_correct)
	values ($1, $2, $3, $4)`

	_, err := s.db.Query(
		query,
		q.SessionID,
		q.QuestionID,
		q.Answer,
		q.IsCorrect,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *responseStore) IsQuizCompleted(sessionID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
        SELECT COUNT(DISTINCT question_id)
        FROM quiz_responses
        WHERE quiz_session_id = $1
    `, sessionID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count >= 5, nil
}

func (s *responseStore) GetQuizResult(sessionID string) (*models.QuizResultResponse, error) {
	var totalQuestions, correctAnswers int
	err := s.db.QueryRow(`
		SELECT
			COUNT(DISTINCT question_id) AS total_questions,
			COUNT(DISTINCT CASE WHEN is_correct THEN question_id END) AS correct_answers
		FROM quiz_responses
		WHERE quiz_session_id = $1;
	`, sessionID).Scan(&totalQuestions, &correctAnswers)

	if err != nil {
		return nil, err
	}

	resp := &models.QuizResultResponse{
		TotalQuestions: totalQuestions,
		CorrectAnswers: correctAnswers,
	}

	return resp, nil
}
