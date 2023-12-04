package store

import (
	"context"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type Message struct {
	db  *pgx.Conn
	log logging.Logger
}

type Messager interface {
	Create(ctx context.Context, message models.Message) (*models.Message, error)
	UpdateStatus(ctx context.Context, id uint64, status models.MessageStatus) error
	GetByID(ctx context.Context, id uint64) (*models.Message, error)
}

func NewMessage(db *pgx.Conn, log logging.Logger) Messager {
	return Message{
		db:  db,
		log: log,
	}
}

func (m Message) Create(ctx context.Context, message models.Message) (*models.Message, error) {
	sql := `
		INSERT INTO messages (status, subject, text, user_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, subject, text, user_id
	`
	result := &models.Message{}
	row := m.db.QueryRow(ctx, sql, message.Status, message.Subject, message.Text, message.UserID)
	err := row.Scan(&result.ID, &result.Status, &result.Subject, &result.Text, &result.UserID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m Message) UpdateStatus(ctx context.Context, id uint64, status models.MessageStatus) error {
	sql := `UPDATE messages SET status = $1 WHERE id = $2;`

	if _, err := m.db.Exec(ctx, sql, status, id); err != nil {
		return err
	}

	return nil
}

func (m Message) GetByID(ctx context.Context, id uint64) (*models.Message, error) {
	sql := `
		SELECT id, status, subject, text, user_id
		FROM messages
		WHERE id = $1;
	`
	message := &models.Message{}
	err := m.db.QueryRow(ctx, sql, id).Scan(&message.ID, &message.Status, &message.Subject, &message.Text, &message.UserID)
	if err != nil {
		return nil, err
	}
	return message, nil
}
