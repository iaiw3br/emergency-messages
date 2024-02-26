package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"projects/emergency-messages/internal/errorx"
)

// Message is a type representing a message
type Message struct {
	ID         uuid.UUID     `json:"id"`
	Subject    string        `json:"subject"`
	Text       string        `json:"text"`
	Status     MessageStatus `json:"status"`
	ReceiverID uuid.UUID     `json:"receiver_id"`
}

// MessageStatus is a type representing a message status
type MessageStatus string

const (
	// Created is a message status when it is created
	Created MessageStatus = "created"
	// Delivered is a message status when it is delivered
	Delivered MessageStatus = "delivered"
	// Failed is a message status when it is failed
	Failed MessageStatus = "failed"
)

// MessageRequest is a type representing a message request.
// It is used to consume messages from the controller.
type MessageRequest struct {
	TemplateID uuid.UUID `json:"template_id"`
	City       string    `json:"city"`
	Strength   string    `json:"strength"`
}

// Validate validates the MessageRequest.
func (m *MessageRequest) Validate() error {
	if m.TemplateID == uuid.Nil {
		return fmt.Errorf("invalid template id: %w", errorx.ErrValidation)
	}
	if m.City == "" {
		return fmt.Errorf("invalid city: %w", errorx.ErrValidation)
	}
	if m.Strength == "" {
		return fmt.Errorf("invalid strength: %w", errorx.ErrValidation)
	}
	return nil
}

// MessageConsumer is a type representing a message consumer.
// It is used to consume messages from the queue broker.
type MessageConsumer struct {
	Subject string        `json:"subject"`
	Text    string        `json:"text"`
	Status  MessageStatus `json:"status"`
	City    string        `json:"city"`
}

// MessageEntity is a type representing a message entity.
// It is used to interact with the database.
type MessageEntity struct {
	bun.BaseModel `bun:"table:messages,alias:m"`
	ID            uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4()"`
	Subject       string        `bun:"subject,notnull"`
	Text          string        `bun:"text,notnull"`
	Status        MessageStatus `bun:"status,notnull"`
	ReceiverID    uuid.UUID     `bun:"receiver_id,notnull"`
	Type          ContactType   `bun:"type,notnull"`
	Value         string        `bun:"value,notnull"`
}

// MessageSend is a type representing a message send.
type MessageSend struct {
	ID         uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4()"`
	Subject    string        `bun:"subject,notnull"`
	Text       string        `bun:"text,notnull"`
	Status     MessageStatus `bun:"status,notnull"`
	ReceiverID uuid.UUID     `bun:"receiver_id,notnull"`
	Type       ContactType   `bun:"type,notnull"`
	Value      string        `bun:"value,notnull"`
}
