package models

type CreateQuestionRequest struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type CreateQuizRequest struct {
	ID          string `json:"id"`
	QuestionIDs []int  `json:"qids"`
}

type UpdateQuizRequest struct {
	ID     string `json:"id"`
	Index  int    `json:"idx"`
	Answer string `json:"answer"`
}
