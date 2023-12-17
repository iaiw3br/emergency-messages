package models

import "github.com/google/uuid"

type Message struct {
	ID      uuid.UUID     `json:"id"`
	Subject string        `json:"subject"`
	Text    string        `json:"text"`
	Status  MessageStatus `json:"status"`
	UserID  uuid.UUID     `json:"user_id"`
}

type MessageStatus string

const (
	Created   MessageStatus = "created"
	Delivered MessageStatus = "delivered"
)

type CreateMessage struct {
	TemplateID uuid.UUID `json:"template_id"`
	City       string    `json:"city"`
	Strength   string    `json:"strength"`
}
