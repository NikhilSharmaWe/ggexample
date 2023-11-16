package store

import (
	"database/sql"

	"sigs.k8s.io/ggexample/models"
)

type ResponseStore interface {
	Create(models.CreateResponseRequest) error
	// GetByID(string) (*models.GetQuizResponse, error)
	// GetRandomQuestionsID(int) (*models.GetQuestionsResponse, error)
	// CheckAnswer(id int, answer string) (bool, error)
	// DeleteByID(id string) error
	// UpdateQuiz(models.UpdateQuizRequest) error
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

// func (s *quizStore) DeleteByID(id string) error {
// 	_, err := s.db.Query("delete from quizes where id = $1", id)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return err
// }

// func (s *quizStore) GetByID(id string) (*models.GetQuizResponse, error) {
// 	rows, err := s.db.Query("select id, questionids, progress from quizes where id = $1", id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		return scanIntoQuiz(rows)
// 	}

// 	return nil, fmt.Errorf("quiz %s not found", id)
// }

// func (s *quizStore) UpdateQuiz(req models.UpdateQuizRequest) error {
// 	query := `
//         update quizes
//         set progress[$2] = $3
//         where id = $1
//     `

// 	_, err := s.db.Exec(query, req.ID, req.Index, req.Answer)

// 	return err
// }

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

// func scanIntoQuiz(rows *sql.Rows) (*models.GetQuizResponse, error) {
// 	q := new(models.GetQuizResponse)

// 	var questionIds pq.Int64Array
// 	var pqArray pq.StringArray

// 	err := rows.Scan(
// 		&q.ID,
// 		&questionIds,
// 		&pqArray,
// 	)

// 	q.QuestionIDs = make([]int, len(questionIds))
// 	for i, id := range questionIds {
// 		q.QuestionIDs[i] = int(id)
// 	}
// 	// fmt.Println(questionIds)

// 	q.Progress = pqArray
// 	// q.QuestionIDs = /
// 	return q, err
// }

// func (s *questionStore) getTotalQuestionCount() (int, error) {
// 	var count int
// 	err := s.db.QueryRow("SELECT count(*) FROM questions").Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }
