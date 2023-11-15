package api

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type QuizManager interface {
	Create() (*models.CreateQuizResponse, error)
	GetByID(string) (*models.GetQuizResponse, error)
	DeleteByID(string) (*models.ResultResponse, error)
	UpdateQuiz(models.UpdateQuizRequest) (*models.ResultResponse, error)
	CheckAnswer(id int, answer string) (bool, error)
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
	qIds, err := q.storeDeps.QuestionStore.GetRandomQuestionIds(5)
	if err != nil {
		return nil, err
	}

	quizId := uuid.NewV4().String()
	cReq := models.CreateQuizRequest{
		ID:          quizId,
		QuestionIDs: qIds,
	}

	if err := q.storeDeps.QuizStore.Create(cReq); err != nil {
		return nil, err
	}

	return &models.CreateQuizResponse{
		ID: quizId,
	}, nil
}

func (q *quizManager) GetByID(id string) (*models.GetQuizResponse, error) {
	resp, err := q.storeDeps.QuizStore.GetByID(id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *quizManager) DeleteByID(id string) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuizStore.DeleteByID(id); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully delete quiz with id: %s", id),
	}, nil
}

func (q *quizManager) UpdateQuiz(req models.UpdateQuizRequest) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuizStore.UpdateQuiz(req); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully updated quiz with id: %s", req.ID),
	}, nil
}

func (q *quizManager) CheckAnswer(id int, answer string) (bool, error) {
	correct, err := q.storeDeps.QuestionStore.CheckAnswer(id, answer)
	if err != nil {
		return false, err
	}

	return correct, nil
}
