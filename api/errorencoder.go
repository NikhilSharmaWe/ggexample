package api

import (
	"context"
	"net/http"
)

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
	errorResponse(w, http.StatusInternalServerError, err.Error())
}
