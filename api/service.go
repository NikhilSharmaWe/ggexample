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
	if err := q.storeDeps.QuestionStore.Create(r); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: "successfully created new question",
	}, nil
}

func (q *questionManger) GetByID(id int) (*models.GetQuestionResponse, error) {
	resp, err := q.storeDeps.QuestionStore.GetByID(id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *questionManger) DeleteByID(id int) (*models.ResultResponse, error) {
	if err := q.storeDeps.QuestionStore.DeleteByID(id); err != nil {
		return nil, err
	}

	return &models.ResultResponse{
		Result: fmt.Sprintf("successfully delete question with %d", id),
	}, nil
}

// type QuizService interface {
// 	GetQuestions() []store.Question
// 	CheckAnswers() []string
// }

// type Application struct {
// 	ListenAddr string
// 	StoreDeps  store.Dependency
// }

// func NewApplication(addr string, s store.Dependency) *Application {
// 	return &Application{
// 		ListenAddr: addr,
// 		StoreDeps:  s,
// 	}
// }

// func (app *Application) HandleCreateQuestion(w http.ResponseWriter, r *http.Request) {
// 	q := new(store.Question)
// 	if err := decodeRequest(w, r, &q); err != nil {
// 		ErrorResponse(w, r, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := app.StoreDeps.QuestionStore.Create(q); err != nil {
// 		ErrorResponse(w, r, http.StatusInternalServerError, "failed to create new question")
// 		return
// 	}

// 	writeEncodedResponse(w, http.StatusOK, map[string]string{"result": "created new question successfully"})
// }

// func (app *Application) HandledDeleteQuestion(w http.ResponseWriter, r *http.Request) {
// 	params := httprouter.ParamsFromContext(r.Context())
// 	id, err := strconv.Atoi(params.ByName("id"))
// 	if err != nil {
// 		ErrorResponse(w, r, http.StatusInternalServerError, "failed to parse data")
// 		return
// 	}

// 	if err := app.StoreDeps.QuestionStore.DeleteByID(id); err != nil {
// 		ErrorResponse(w, r, http.StatusInternalServerError, "failed to delete question")
// 		return
// 	}

// 	writeEncodedResponse(w, http.StatusOK, map[string]string{"result": "deleted question successfully"})
// }

// func (app *Application) HandleGetQuestion(w http.ResponseWriter, r *http.Request) {
// 	params := httprouter.ParamsFromContext(r.Context())
// 	id, err := strconv.Atoi(params.ByName("id"))
// 	if err != nil {
// 		ErrorResponse(w, r, http.StatusInternalServerError, "failed to parse data")
// 		return
// 	}

// 	question, err := app.StoreDeps.QuestionStore.GetByID(id)
// 	if err != nil {
// 		ErrorResponse(w, r, http.StatusInternalServerError, "failed to get question")
// 		return
// 	}

// 	writeEncodedResponse(w, http.StatusOK, *question)
// }

// func NewQuizService() *QuizService{
// 	return &quizService{}
// }

// func handleQuestion(w http.ResponseWriter, r *http.Request) {
// 	id, question := getQuestion()
// 	resp := map[string]any{
// 		"question": question.Text,
// 		"options":  question.Options,
// 		"id":       id,
// 	}
// 	writeJSON(w, http.StatusOK, resp)
// }

// func handleAnswerCheck(w http.ResponseWriter, r *http.Request) {
// 	params := httprouter.ParamsFromContext(r.Context())
// 	id := params.ByName("id")

// 	question, ok := questions[id]
// 	if !ok {
// 		errorResponse(w, r, http.StatusBadRequest, InvalidQuestionReqErr)
// 		return
// 	}

// 	resp := make(map[string]any)
// 	if err := readJSON(w, r, &resp); err != nil {
// 		errorResponse(w, r, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	option, ok := resp["answer"]
// 	if !ok {
// 		errorResponse(w, r, http.StatusBadRequest, InvalidAnswerReqErr)
// 		return
// 	}

// 	if option == question.Answer {
// 		writeJSON(w, http.StatusOK, map[string]string{"result": "You are Correct"})
// 	} else {
// 		writeJSON(w, http.StatusOK, map[string]string{"result": "Better luck next time"})
// 	}
// }
