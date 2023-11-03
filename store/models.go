package store

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

func NewQuestion(text string, options []string, answer string) *Question {
	return &Question{
		Text:    text,
		Options: options,
		Answer:  answer,
	}
}
