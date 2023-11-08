package api

import (
	"fmt"

	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type QuestionManager interface {
	Create(models.CreateQuestionRequest) (*models.ResultResponse, error)
	GetByID(int) (*models.GetQuestionResponse, error)
	DeleteByID(int) (*models.ResultResponse, error)
}

type questionManger struct {
	storeDeps store.Dependency
}

func NewQuestionManager(storeDeps store.Dependency) QuestionManager {
	return &questionManger{
		storeDeps: storeDeps,
	}
}

func (q *questionManger) Create(r models.CreateQuestionRequest) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuestionStore.CreateQuestion(r); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: "successfully created new question",
	}, nil
}

func (q *questionManger) GetByID(id int) (*models.GetQuestionResponse, error) {
	resp, err := q.storeDeps.QuestionStore.GetQuestionByID(id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *questionManger) DeleteByID(id int) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuestionStore.DeleteQuestionByID(id); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully delete question with %d", id),
	}, nil
}
