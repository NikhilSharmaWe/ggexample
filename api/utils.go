package api

import (
	"encoding/json"
	"net/http"
)

// type Question struct {
// 	Text    string   `json:"question"`
// 	Options []string `json:"options"`
// 	Answer  string   `json:"answer"`
// }

// func NewQuestion(text string, options []string, answer string) *Question {
// 	return &Question{
// 		Text:    text,
// 		Options: options,
// 		Answer:  answer,
// 	}
// }

// var Questions = map[string]Question{
// 	"0": {
// 		Text:    "Which is the best Rock Band",
// 		Options: []string{"Pink Floyd", "Radio Head", "Audio Slave"},
// 		Answer:  "Audio Slave",
// 	},
// 	"1": {
// 		Text:    "Who is the Director of Oppenhiemer",
// 		Options: []string{"Tarintino", "Nolan", "Nikhil"},
// 		Answer:  "Nolan",
// 	},
// 	"2": {
// 		Text:    "Who is the MMA Goat",
// 		Options: []string{"CM Punk", "Paddy", "Volkanovski"},
// 		Answer:  "Volkanovski",
// 	},
// }

// func getQuestion() (string, Question) {
// 	id := strconv.Itoa(rand.Intn(len(questions)))
// 	return id, questions[id]
// }

func writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
