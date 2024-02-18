// Package postgres
// Description: This file contains the implementation of the message store interface for the postgres database.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"projects/emergency-messages/internal/models"
)

type MessageStore struct {
	db *bun.DB
}

func NewMessage(db *bun.DB) *MessageStore {
	return &MessageStore{
		db: db,
	}
}

// Create creates the struct of a message in the database.
// It takes in a context, the new struct of the message.
// It returns an error if the create operation fails.
func (s *MessageStore) Create(ctx context.Context, m *models.MessageEntity) error {
	_, err := s.db.
		NewInsert().
		Model(m).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("creating message: couldn't create: %v. Error: %w", m, err)
	}
	return nil
}

// UpdateStatus updates the status of a message in the database.
// It takes in a context, the ID of the message, and the new status.
// It returns an error if the update operation fails.
func (s *MessageStore) UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error {
	exec, err := s.db.
		NewUpdate().
		Model(&models.MessageEntity{}).
		Set("status = ?", string(status)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("updating message: couldn't update with id: %s. Error: %w", id, err)
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return fmt.Errorf("updating message: couldn't get the number of rows affected with id: %s. Error: %w", id, err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// FindByStatus retrieves messages from the database by status.
// It takes in a context and the status of the messages.
// It returns a slice of message entities and an error if the find operation fails.
func (s *MessageStore) FindByStatus(ctx context.Context, status models.MessageStatus) ([]models.MessageEntity, error) {
	entities := make([]models.MessageEntity, 0)

	err := s.db.
		NewSelect().
		Model(&entities).
		Where("status = ?", string(status)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding messages by status: couldn't find messages by status: %s. Error: %w", status, err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	return entities, nil
}
