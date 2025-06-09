package models

type Poll struct {
	ID            string   `json:"id"`
	Question      string   `json:"question"`
	AnswerOptions []Option `json:"options"` 
	IsClosed bool `json:"closed"` 
}

type Option struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Votes   int    `json:"votes"`
}
