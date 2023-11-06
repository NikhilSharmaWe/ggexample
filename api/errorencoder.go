package api

import (
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := map[string]any{
		"error": message,
	}

	err := writeEncodedResponse(w, status, env)
	if err != nil {
		w.WriteHeader(status)
	}
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusNotFound, NotFoundErr)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusMethodNotAllowed, MethodNotFoundErr)
}
