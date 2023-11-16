package api

import (
	"fmt"

	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type ResponseManager interface {
	Create(models.CreateResponseRequest) (*models.ResultResponse, error)
}

type responseManager struct {
	storeDeps store.Dependency
}

func NewResponseManager(storeDeps store.Dependency) ResponseManager {
	return &responseManager{
		storeDeps: storeDeps,
	}
}

func (q *responseManager) Create(req models.CreateResponseRequest) (*models.ResultResponse, error) {
	correct, err := q.storeDeps.QuestionStore.CheckAnswer(req.QuestionID, req.Answer)

	if err != nil {
		return nil, err
	}

	req.IsCorrect = correct

	if err := q.storeDeps.ResponseStore.Create(req); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully stored response for quiz session with id: %s", req.SessionID),
	}, nil
}
