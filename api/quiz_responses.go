package api

import (
	"fmt"

	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type ResponseManager interface {
	Create(models.CreateResponseRequest) (*models.ResultResponse, error)
	// GetByID(string) (*models.GetQuizResponse, error)
	// DeleteByID(string) (*models.ResultResponse, error)
	// UpdateQuiz(models.UpdateQuizRequest) (*models.ResultResponse, error)
	// CheckAnswer(id int, answer string) (bool, error)
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

// func (q *quizManager) GetByID(id string) (*models.GetQuizResponse, error) {
// 	resp, err := q.storeDeps.QuizStore.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (q *quizManager) DeleteByID(id string) (*models.ResultResponse, error) {
// 	if err := q.storeDeps.QuizStore.DeleteByID(id); err != nil {
// 		return nil, err
// 	}

// 	return &models.ResultResponse{
// 		Result: fmt.Sprintf("successfully delete quiz with id: %s", id),
// 	}, nil
// }

// func (q *quizManager) UpdateQuiz(req models.UpdateQuizRequest) (*models.ResultResponse, error) {
// 	if err := q.storeDeps.QuizStore.UpdateQuiz(req); err != nil {
// 		return nil, err
// 	}

// 	return &models.ResultResponse{
// 		Result: fmt.Sprintf("successfully updated quiz with id: %s", req.ID),
// 	}, nil
// }

// func (q *quizManager) CheckAnswer(id int, answer string) (bool, error) {
// 	correct, err := q.storeDeps.QuestionStore.CheckAnswer(id, answer)
// 	if err != nil {
// 		return false, err
// 	}

// 	return correct, nil
// }
