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
	"sigs.k8s.io/ggexample/models"
)

// var (
// 	upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}
// )

// type connectionResponse struct {
// 	request *http.Request
// 	svc     QuizManager
// }

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
		id, err := strconv.Atoi(r.(string))
		if err != nil {
			return nil, Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
		}
		return svc.GetByID(id)
	}
}

func makeCreateQuizEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		resp, err := svc.Create()
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func makeGetQuizEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		id := r.(string)
		resp, err := svc.GetByID(id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func makeGetQuizQuestionEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		id := r.(string)
		completed, err := svc.IsQuizCompleted(id)
		if err != nil {
			return nil, err
		}

		if completed {
			resp, err := svc.GetQuizResult(id)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		resp, err := svc.Get(id)
		if err != nil {
			return nil, err
		}

		return &models.Question{
			ID:      resp.ID,
			Text:    resp.Text,
			Options: resp.Options,
		}, nil
	}
}

func makeDeleteQuizEndpoint(svc QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		id := r.(string)
		resp, err := svc.DeleteByID(id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func makeResponseEndpoint(responseSVC ResponseManager, quizSVC QuizManager) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(models.CreateResponseRequest)

		id := req.SessionID
		completed, err := quizSVC.IsQuizCompleted(id)
		if err != nil {
			return nil, err
		}

		if completed {
			resp, err := quizSVC.GetQuizResult(id)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		resp, err := responseSVC.Create(req)

		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// func makeCheckQuiz(svc QuizManager) endpoint.Endpoint {
// 	return func(ctx context.Context, r interface{}) (interface{}, error) {
// 		req := r.(*http.Request)
// 		return returnConnectionResponse(req, svc)
// 	}
// }

// func returnConnectionResponse(r *http.Request, svc QuizManager) (any, error) {
// 	return connectionResponse{
// 		request: r,
// 		svc:     svc,
// 	}, nil
// }

func NewHTTPHandler(questionSVC QuestionManager, quizSVC QuizManager, responseSVC ResponseManager) http.Handler {

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

	createQuizHandler := httptransport.NewServer(
		makeCreateQuizEndpoint(quizSVC),
		decodeEmptyRequest,
		encodeResponse,
		options...,
	)

	getQuizHandler := httptransport.NewServer(
		makeGetQuizQuestionEndpoint(quizSVC),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	deleteQuizHandler := httptransport.NewServer(
		makeDeleteQuizEndpoint(quizSVC),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	responseHandler := httptransport.NewServer(
		makeResponseEndpoint(responseSVC, quizSVC),
		decodeResponseRequest,
		encodeResponse,
		options...,
	)

	// checkQuizHandler := httptransport.NewServer(
	// 	makeCheckQuiz(quizSVC),
	// 	returnRequest,
	// 	handleWebsocketResponse,
	// 	options...,
	// )

	mux.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	mux.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	mux.Handle("/question/create", createQuestionHandler).Methods(http.MethodPost)
	mux.Handle("/question/delete/{id}", deleteQuestionHandler).Methods(http.MethodDelete)
	mux.Handle("/question/get/{id}", getQuestionHandler).Methods(http.MethodGet)

	mux.Handle("/quiz/create", createQuizHandler).Methods(http.MethodGet)
	mux.Handle("/quiz/get/{id}", getQuizHandler).Methods(http.MethodGet)
	mux.Handle("/quiz/delete/{id}", deleteQuizHandler).Methods(http.MethodDelete)
	mux.Handle("/quiz/response", responseHandler).Methods(http.MethodPost)

	// mux.Handle("/websocket", checkQuizHandler).Methods(http.MethodGet)

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
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, Error{
			Message: "bad request",
			Code:    http.StatusBadRequest,
		}
	}
	return id, nil
}

func decodeResponseRequest(_ context.Context, r *http.Request) (any, error) {
	req := models.CreateResponseRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, Error{
			Message: err.Error(),
			Code:    http.StatusBadGateway,
		}
	}

	return req, nil

}

func decodeEmptyRequest(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func returnRequest(_ context.Context, r *http.Request) (any, error) {
	return r, nil
}

// func handleWebsocketResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
// 	r := response.(connectionResponse)
// 	ws, err := upgrader.Upgrade(w, r.request, nil)
// 	if err != nil {
// 		return err
// 	}

// 	var submitted bool

// 	for !submitted {
// 		message := models.WebsocketMessage{
// 			Data: models.QuestionAnswerResponse{},
// 		}
// 		if err := ws.ReadJSON(&message); err != nil {
// 			return err
// 		}

// 		switch message.Event {
// 		case "answer":
// 			data := message.Data.(map[string]interface{})
// 			id, err := strconv.Atoi(data["id"].(string))
// 			if err != nil {
// 				return err
// 			}
// 			answer := data["answer"].(string)

// 			correct, err := r.svc.CheckAnswer(id, answer)
// 			if err != nil {
// 				return err
// 			}

// 			var content string
// 			if correct {
// 				content = "correct"
// 			} else {
// 				content = "wrong"
// 			}

// 			if err := ws.WriteJSON(&models.WebsocketMessage{
// 				Event: "answer",
// 				Data: models.QuestionResultResponse{
// 					Index:  data["idx"].(string),
// 					Result: content,
// 				},
// 			}); err != nil {
// 				return err
// 			}

// 		case "submit":
// 			resp := make(map[int]string)

// 			for idx, res := range message.Data.(map[string]interface{}) {
// 				res := res.(map[string]interface{})
// 				idx, err := strconv.Atoi(idx)
// 				if err != nil {
// 					return err
// 				}

// 				id, err := strconv.Atoi(res["id"].(string))
// 				if err != nil {
// 					return err
// 				}

// 				correct, err := r.svc.CheckAnswer(id, res["answer"].(string))
// 				if err != nil {
// 					return err
// 				}

// 				if correct {
// 					resp[idx+1] = "correct"
// 				} else {
// 					resp[idx+1] = "wrong"
// 				}
// 			}

// 			if err := ws.WriteJSON(&models.WebsocketMessage{
// 				Event: "result",
// 				Data:  resp,
// 			}); err != nil {
// 				return err
// 			}

// 			submitted = true
// 		}
// 	}

// 	return nil
// }

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
