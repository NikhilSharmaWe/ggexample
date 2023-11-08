package models

type GetQuestionResponse struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type GetQuestionsResponse struct {
	Questions []*GetQuestionResponse `json:"questions"`
}

type QuestionAnswerResponse struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

type ResultResponse struct {
	Result string `json:"result"`
}
