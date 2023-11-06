package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"sigs.k8s.io/ggexample/store"
)

type Application struct {
	listenAddr string
	storeDeps  store.Dependency
}

func NewApplication(addr string, s store.Dependency) *Application {
	return &Application{
		listenAddr: addr,
		storeDeps:  s,
	}
}

func (app *Application) Start() {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)

	router.HandlerFunc(http.MethodPost, "/createQ", app.handleCreateQuestion)
	router.HandlerFunc(http.MethodDelete, "/deleteQ/:id", app.handledDeleteQuestion)
	router.HandlerFunc(http.MethodGet, "/getQ/:id", app.handleGetQuestion)

	// router.HandlerFunc(http.MethodGet, "/question", handleQuestion)
	// router.HandlerFunc(http.MethodPost, "/check/:id", handleAnswerCheck)

	log.Printf("Starting server at %s", app.listenAddr)
	log.Fatal(http.ListenAndServe(app.listenAddr, router))
}
