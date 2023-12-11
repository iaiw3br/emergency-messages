package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/service"
	"github.com/uptrace/bun"
)

type messageStore struct {
	db *bun.DB
}

type messageEntity struct {
	bun.Model
	ID      string        `json:"id"`
	Subject string        `json:"subject"`
	Text    string        `json:"text"`
	Status  MessageStatus `json:"status"`
	UserID  uint64        `json:"user_id"`
}

type MessageStatus string

func NewMessage(db *bun.DB) service.Messager {
	return &messageStore{
		db: db,
	}
}

func (s *messageStore) Create(ctx context.Context, m *models.Message) error {
	entity := messageEntity{}
	_, err := s.db.
		NewInsert().
		Model(messageEntity{}).
		Exec(ctx, m)

	if err != nil {
		return err
	}
	m.ID = entity.ID
	return nil
}

func (s *messageStore) UpdateStatus(ctx context.Context, id string, status models.MessageStatus) error {
	exec, err := s.db.
		NewUpdate().
		Model(messageEntity{}).
		Set("status = ?", status).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New(fmt.Sprintf("message was not found by id: %s", id))
	}
	return nil
}
