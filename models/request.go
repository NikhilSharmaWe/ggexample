package models

type CreateQuestionRequest struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type CreateQuizRequest struct {
	ID string `json:"id"`
	// QuestionIDs []int  `json:"qids"`
}

type CreateResponseRequest struct {
	SessionID  string `json:"session_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
}

type UpdateQuizRequest struct {
	ID     string `json:"id"`
	Index  int    `json:"idx"`
	Answer string `json:"answer"`
}
