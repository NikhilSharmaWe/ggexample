package api

import (
	"sigs.k8s.io/ggexample/models"
	"sigs.k8s.io/ggexample/store"
)

type QuizManager interface {
	GetQuiz() (*models.GetQuestionsResponse, error)
}

type quizManager struct {
	storeDeps store.Dependency
}

func NewQuizManager(storeDeps store.Dependency) QuizManager {
	return &quizManager{
		storeDeps: storeDeps,
	}
}

func (q *quizManager) GetQuiz() (*models.GetQuestionsResponse, error) {
	resp, err := q.storeDeps.QuestionStore.GetRandomQuestions(5)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
