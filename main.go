package main

import (
	"log"
	"net/http"

	"sigs.k8s.io/ggexample/api"
	"sigs.k8s.io/ggexample/store"
)

func main() {
	db, err := store.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	questionStore := store.NewQuestionStore(db)
	if err := questionStore.InitQuestionRelation(); err != nil {
		log.Fatal(err)
	}

	quizStore := store.NewQuizStore(db)
	if err := quizStore.InitQuizRelation(); err != nil {
		log.Fatal(err)
	}

	responseStore := store.NewResponseStore(db)
	if err := responseStore.InitResponseRelation(); err != nil {
		log.Fatal(err)
	}

	storeDep := store.Dependency{
		QuestionStore: questionStore,
		QuizStore:     quizStore,
		ResponseStore: responseStore,
	}

	questionSVC := api.NewQuestionManager(storeDep)
	quizSVC := api.NewQuizManager(storeDep)
	responseSVC := api.NewResponseManager(storeDep)

	router := api.NewHTTPHandler(questionSVC, quizSVC, responseSVC)

	log.Printf("Starting server at %s", ":8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
