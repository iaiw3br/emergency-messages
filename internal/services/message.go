package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
)

type MessageService struct {
	producer      Producer
	templateStore Template
	log           *slog.Logger
}

type Producer interface {
	Send(messageBytes []byte) error
}

type Template interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateEntity, error)
}

func NewMessage(producer Producer, templateStore Template, log *slog.Logger) *MessageService {
	return &MessageService{
		producer:      producer,
		templateStore: templateStore,
		log:           log,
	}
}

func (s *MessageService) Send(ctx context.Context, message models.MessageRequest) error {
	// validate message
	if err := message.Validate(); err != nil {
		return err
	}

	// get template by id
	template, err := s.templateStore.GetByID(ctx, message.TemplateID)
	if err != nil {
		s.log.With(slog.Any("templateID", message.TemplateID)).
			Error("getting template", err)
		if err == sql.ErrNoRows {
			return errorx.ErrNotFound
		}
		return errorx.ErrInternal
	}

	newMessage := models.MessageConsumer{
		Subject: template.Subject,
		Text:    fmt.Sprintf(template.Text, message.City, message.Strength),
		Status:  models.Created,
		City:    message.City,
	}

	messageBytes, err := json.Marshal(newMessage)
	if err != nil {
		s.log.With(slog.Any("message", newMessage)).
			Error("marshaling message", err)
		return errorx.ErrInternal
	}

	// send to queue
	if err = s.producer.Send(messageBytes); err != nil {
		s.log.With(slog.Any("message", newMessage)).
			Error("sending message", err)
		return errorx.ErrInternal

	}

	return nil
}
