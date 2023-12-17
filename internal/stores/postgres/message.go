package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"

	"github.com/uptrace/bun"
)

type messageStore struct {
	db *bun.DB
}

func NewMessage(db *bun.DB) services.Message {
	return &messageStore{
		db: db,
	}
}

// Create creates the struct of a message in the database.
// It takes in a context, the new struct of the message.
// It returns an error if the create operation fails.
func (s *messageStore) Create(ctx context.Context, m *models.MessageEntity) error {
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
func (s *messageStore) UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error {
	exec, err := s.db.
		NewUpdate().
		Model(&models.TemplateEntity{}).
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
		return fmt.Errorf("updating message: couldn't find message with id: %s", id)
	}
	return nil
}
