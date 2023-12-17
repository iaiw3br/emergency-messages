package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

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

type MessageEntity struct {
	bun.BaseModel `bun:"table:messages,alias:m"`
	ID            uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4()"`
	Subject       string        `bun:"subject,notnull"`
	Text          string        `bun:"text,notnull"`
	Status        MessageStatus `bun:"status,notnull"`
	UserID        uuid.UUID     `bun:"user_id,notnull"`
}
