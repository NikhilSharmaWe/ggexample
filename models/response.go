package models

type GetQuestionResponse struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

type ResultResponse struct {
	Result string `json:"result"`
}
