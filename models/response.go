package models

type GetQuestionResponse struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type QuizResultResponse struct {
	TotalQuestions int `json:"total_questions"`
	CorrectAnswers int `json:"correct_answers"`
}

type CreateQuizResponse struct {
	ID string `json:"id"`
}

type ResultResponse struct {
	Result string `json:"result"`
}
