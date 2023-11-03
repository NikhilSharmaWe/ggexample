package main

import (
	"log"

	"sigs.k8s.io/ggexample/api"
	"sigs.k8s.io/ggexample/store"
)

func main() {
	s, err := store.New()
	if err != nil {
		log.Fatal(err)
	}
	s.Init()

	app := api.NewApplication(":8000", s)
	app.Start()

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
