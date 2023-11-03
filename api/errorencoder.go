package api

import (
	"net/http"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := map[string]any{
		"error": message,
	}

	err := writeJSON(w, status, env)
	if err != nil {
		w.WriteHeader(status)
	}
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, NotFoundErr)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusMethodNotAllowed, MethodNotFoundErr)
}
