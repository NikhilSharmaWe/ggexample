package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"sigs.k8s.io/ggexample/models"
)

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

func NewHTTPHandler(svc QuestionManager) http.Handler {

	mux := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	createQuestionHandler := httptransport.NewServer(
		makeCreateQuestionEndpoint(svc),
		decodeCreateQuestionRequest,
		encodeResponse,
		options...,
	)

	deleteQuestionHandler := httptransport.NewServer(
		makeDeleteQuestionEndpoint(svc),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	getQuestionHandler := httptransport.NewServer(
		makeGetQuestionEndpoint(svc),
		decodeParamIDRequest,
		encodeResponse,
		options...,
	)

	mux.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	mux.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	mux.Handle("/createQ", createQuestionHandler).Methods(http.MethodPost)
	mux.Handle("/deleteQ/{id}", deleteQuestionHandler).Methods(http.MethodDelete)
	mux.Handle("/getQ/{id}", getQuestionHandler).Methods(http.MethodGet)

	return mux
}

func decodeCreateQuestionRequest(_ context.Context, r *http.Request) (any, error) {
	req := models.CreateQuestionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeParamIDRequest(_ context.Context, r *http.Request) (any, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		return nil, err
	}

	return id, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
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
