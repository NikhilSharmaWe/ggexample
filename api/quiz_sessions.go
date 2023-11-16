package api

import (
	"fmt"

	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type QuizManager interface {
	Create() (*models.CreateQuizResponse, error)
	Get(sessionId string) (*models.GetQuestionResponse, error)
	Exists(sessionId string) (bool, error)
	DeleteByID(string) (*models.ResultResponse, error)
	CheckAnswer(id int, answer string) (bool, error)
	IsQuizCompleted(id string) (bool, error)
	GetQuizResult(sessionID string) (*models.QuizResultResponse, error)
}

type quizManager struct {
	storeDeps store.Dependency
}

func NewQuizManager(storeDeps store.Dependency) QuizManager {
	return &quizManager{
		storeDeps: storeDeps,
	}
}

func (q *quizManager) Create() (*models.CreateQuizResponse, error) {
	resp, err := q.storeDeps.QuizStore.Create()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *quizManager) Get(sessionId string) (*models.GetQuestionResponse, error) {
	resp, err := q.storeDeps.QuestionStore.GetNextQuestion(sessionId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *quizManager) Exists(sessionId string) (bool, error) {
	resp, err := q.storeDeps.QuizStore.Exists(sessionId)
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (q *quizManager) IsQuizCompleted(sessionID string) (bool, error) {
	resp, err := q.storeDeps.ResponseStore.IsQuizCompleted(sessionID)
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (q *quizManager) GetQuizResult(sessionID string) (*models.QuizResultResponse, error) {
	resp, err := q.storeDeps.ResponseStore.GetQuizResult(sessionID)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (q *quizManager) DeleteByID(id string) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuizStore.DeleteByID(id); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully delete quiz with id: %s", id),
	}, nil
}

func (q *quizManager) CheckAnswer(id int, answer string) (bool, error) {
	correct, err := q.storeDeps.QuestionStore.CheckAnswer(id, answer)
	if err != nil {
		return false, err
	}

	return correct, nil
}
