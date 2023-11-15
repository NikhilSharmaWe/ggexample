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

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type QuestionAnswerResponse struct {
	ID     string `json:"id"`
	Index  int    `json:"idx"`
	Answer string `json:"answer"`
}

type QuestionResultResponse struct {
	Index  string `json:"idx"`
	Result string `json:"result"`
}

type CreateQuizResponse struct {
	ID string `json:"id"`
}

type GetQuizResponse struct {
	ID          string   `json:"id"`
	QuestionIDs []int    `json:"qids"`
	Progress    []string `json:"progress"`
}

type ResultResponse struct {
	Result string `json:"result"`
}
