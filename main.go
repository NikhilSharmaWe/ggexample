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

	questionSVC := api.NewQuestionManager(store.Dependency{
		QuestionStore: questionStore,
		QuizStore:     quizStore,
	})

	quizSVC := api.NewQuizManager(store.Dependency{
		QuestionStore: questionStore,
		QuizStore:     quizStore,
	})

	router := api.NewHTTPHandler(questionSVC, quizSVC)

	log.Printf("Starting server at %s", ":8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
