package models

type Message struct {
	ID      string        `json:"id"`
	Subject string        `json:"subject"`
	Text    string        `json:"text"`
	Status  MessageStatus `json:"status"`
	UserID  string        `json:"user_id"`
}

type MessageStatus string

const (
	Created   MessageStatus = "created"
	Delivered MessageStatus = "delivered"
)

type CreateMessage struct {
	TemplateID string `json:"template_id"`
	City       string `json:"city"`
	Strength   string `json:"strength"`
}
