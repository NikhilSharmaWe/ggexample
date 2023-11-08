package api

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/ggexample/models"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Application struct {
	QuestionManager QuestionManager
	QuizManager     QuizManager
}

func NewApplication(questionSVC QuestionManager, quizSVC QuizManager) *Application {
	return &Application{
		QuestionManager: questionSVC,
		QuizManager:     quizSVC,
	}
}

func makeCreateQuestionEndpoint(svc QuestionManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(models.CreateQuestionRequest)
		return svc.Create(req)
	}
}

func makeDeleteQuestionEndpoint(svc QuestionManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(int)
		return svc.DeleteByID(req)
	}
}

func makeGetQuestionEndpoint(svc QuestionManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(int)
		return svc.GetByID(req)
	}
}

func makeGetQuizEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return svc.GetQuiz()
	}
}

func (app *Application) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		writeEncodedResponse(w, http.StatusInternalServerError, err)
		return
	}

	var result map[int]models.QuestionAnswerResponse

	if err := ws.ReadJSON(&result); err != nil {
		writeEncodedResponse(w, http.StatusInternalServerError, err)
		return
	}

	resp := make(map[int]string)

	for idx, res := range result {
		id, err := strconv.Atoi(res.ID)
		if err != nil {
			writeEncodedResponse(w, http.StatusInternalServerError, err)
			return
		}
		correct, err := app.QuestionManager.CheckAnswer(id, res.Answer)
		if err != nil {
			writeEncodedResponse(w, http.StatusInternalServerError, err)
			return
		}

		if correct {
			resp[idx+1] = "correct"
		} else {
			resp[idx+1] = "wrong"
		}
	}

	if err := ws.WriteJSON(resp); err != nil {
		writeEncodedResponse(w, http.StatusInternalServerError, err)
		return
	}
}

func NewHTTPHandler(questionSVC QuestionManager, quizSVC QuizManager) http.Handler {

	mux := mux.NewRouter()
	app := NewApplication(questionSVC, quizSVC)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	createQuestionHandler := httptransport.NewServer(
		makeCreateQuestionEndpoint(questionSVC),
		decodeCreateQuestionRequest,
		encodeResponse,
		options...,
	)

	deleteQuestionHandler := httptransport.NewServer(
		makeDeleteQuestionEndpoint(questionSVC),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	getQuestionHandler := httptransport.NewServer(
		makeGetQuestionEndpoint(questionSVC),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	getQuizHandler := httptransport.NewServer(
		makeGetQuizEndpoint(quizSVC),
		decodeEmptyRequest,
		renderQuiz,
		options...,
	)

	mux.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	mux.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	mux.Handle("/question/create", createQuestionHandler).Methods(http.MethodPost)
	mux.Handle("/question/delete/{id}", deleteQuestionHandler).Methods(http.MethodDelete)
	mux.Handle("/question/get/{id}", getQuestionHandler).Methods(http.MethodGet)

	mux.Handle("/quiz", getQuizHandler).Methods(http.MethodGet)

	mux.HandleFunc("/websocket", app.HandleConnections)

	return mux
}

func decodeCreateQuestionRequest(_ context.Context, r *http.Request) (any, error) {
	req := models.CreateQuestionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, Error{
			Message: err.Error(),
			Code:    http.StatusBadGateway,
		}
	}

	if req.Text == "" || len(req.Options) != 3 || req.Answer == "" {
		return nil, Error{
			Message: "invalid data",
			Code:    http.StatusBadRequest,
		}
	}

	return req, nil
}

func decodeParamIDRequest(_ context.Context, r *http.Request) (any, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		return nil, Error{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	return id, nil
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func renderQuiz(_ context.Context, w http.ResponseWriter, response interface{}) error {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.Execute(w, response); err != nil {
		log.Println(err)
	}

	return nil
}

func writeEncodedResponse(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
