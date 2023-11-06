package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"sigs.k8s.io/ggexample/api"
	"sigs.k8s.io/ggexample/store"
)

func main() {
	s, err := store.NewQuestionStore()
	if err != nil {
		log.Fatal(err)
	}
	s.InitQuestionRelation()

	app := api.NewApplication(":8000", store.Dependency{
		QuestionStore: s,
	})

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(api.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowed)

	router.HandlerFunc(http.MethodPost, "/createQ", app.HandleCreateQuestion)
	router.HandlerFunc(http.MethodDelete, "/deleteQ/:id", app.HandledDeleteQuestion)
	router.HandlerFunc(http.MethodGet, "/getQ/:id", app.HandleGetQuestion)

	// router.HandlerFunc(http.MethodGet, "/question", handleQuestion)
	// router.HandlerFunc(http.MethodPost, "/check/:id", handleAnswerCheck)

	log.Printf("Starting server at %s", app.ListenAddr)
	log.Fatal(http.ListenAndServe(app.ListenAddr, router))

	// api.StartService()

	// if err = s.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// question := api.Questions["0"]
	// q := store.NewQuestion(question.Text, question.Options, question.Answer)
	// if err := s.CreateQuestion(q); err != nil {
	// 	log.Fatal(err)
	// }

	// s.DeleteQuestion(1)

	// q, err = s.GetQuestionByID(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%+v\n", q)

}
