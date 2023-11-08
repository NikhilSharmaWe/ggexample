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

type connectionResponse struct {
	request *http.Request
	svc     QuizManager
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

func makeCheckQuiz(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(*http.Request)
		return returnConnectionResponse(req, svc)
	}
}

func makeGetQuizEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return svc.GetQuiz()
	}
}

func returnConnectionResponse(r *http.Request, svc QuizManager) (any, error) {
	return connectionResponse{
		request: r,
		svc:     svc,
	}, nil
}

func NewHTTPHandler(questionSVC QuestionManager, quizSVC QuizManager) http.Handler {

	mux := mux.NewRouter()

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

	checkQuizHandler := httptransport.NewServer(
		makeCheckQuiz(quizSVC),
		returnRequest,
		handleWebsocketResponse,
		options...,
	)

	mux.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	mux.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	mux.Handle("/question/create", createQuestionHandler).Methods(http.MethodPost)
	mux.Handle("/question/delete/{id}", deleteQuestionHandler).Methods(http.MethodDelete)
	mux.Handle("/question/get/{id}", getQuestionHandler).Methods(http.MethodGet)

	mux.Handle("/quiz", getQuizHandler).Methods(http.MethodGet)
	mux.Handle("/websocket", checkQuizHandler).Methods(http.MethodGet)

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

func returnRequest(_ context.Context, r *http.Request) (any, error) {
	return r, nil
}

type AnswerData struct {
	ID     int    `json:"id"`
	Idx    int    `json:"idx"`
	Answer string `json:"answer"`
}

func handleWebsocketResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(connectionResponse)
	ws, err := upgrader.Upgrade(w, r.request, nil)
	if err != nil {
		return err
	}

	var submitted bool

	for !submitted {
		message := models.WebsocketMessage{
			Data: models.QuestionAnswerResponse{},
		}
		if err := ws.ReadJSON(&message); err != nil {
			return err
		}

		switch message.Event {
		case "answer":
			data := message.Data.(map[string]interface{})
			id, err := strconv.Atoi(data["id"].(string))
			if err != nil {
				return err
			}
			answer := data["answer"].(string)

			correct, err := r.svc.CheckAnswer(id, answer)
			if err != nil {
				return err
			}

			var content string
			if correct {
				content = "correct"
			} else {
				content = "wrong"
			}

			if err := ws.WriteJSON(&models.WebsocketMessage{
				Event: "answer",
				Data: models.QuestionResult{
					Index:  data["idx"].(string),
					Result: content,
				},
			}); err != nil {
				return err
			}

		case "submit":
			resp := make(map[int]string)

			for idx, res := range message.Data.(map[string]interface{}) {
				res := res.(map[string]interface{})
				idx, err := strconv.Atoi(idx)
				if err != nil {
					return err
				}

				id, err := strconv.Atoi(res["id"].(string))
				if err != nil {
					return err
				}

				correct, err := r.svc.CheckAnswer(id, res["answer"].(string))
				if err != nil {
					return err
				}

				if correct {
					resp[idx+1] = "correct"
				} else {
					resp[idx+1] = "wrong"
				}
			}

			if err := ws.WriteJSON(&models.WebsocketMessage{
				Event: "result",
				Data:  resp,
			}); err != nil {
				return err
			}

			submitted = true
		}
	}

	return nil
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
