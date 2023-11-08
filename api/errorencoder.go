package api

import (
	"context"
	"net/http"
)

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

func errorResponse(w http.ResponseWriter, status int, message any) {
	resp := map[string]any{
		"error": message,
	}
	err := writeEncodedResponse(w, status, resp)
	if err != nil {
		w.WriteHeader(status)
	}
}

func notFoundResponse(w http.ResponseWriter, _ *http.Request) {
	errorResponse(w, http.StatusNotFound, NotFoundErr)
}

func methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	errorResponse(w, http.StatusMethodNotAllowed, MethodNotFoundErr)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	e, ok := err.(Error)
	if !ok {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	errorResponse(w, e.Code, e.Message)
}
