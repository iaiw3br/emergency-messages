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

func NewMessage(db *pgx.Conn, log logging.Logger) Message {
	return Message{
		db:  db,
		log: log,
	}
}

func (m Message) Create(ctx context.Context, message models.Message) error {
	sql := `
		INSERT INTO messages (status, subject, text, user_id) 
		VALUES ($1, $2, $3, $4);
	`
	_, err := m.db.Exec(ctx, sql, message.Status, message.Subject, message.Text, message.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (m Message) UpdateStatus(ctx context.Context, id uint64, status models.MessageStatus) error {
	sql := `UPDATE messages SET status = $1 WHERE id = $2;`

	if _, err := m.db.Exec(ctx, sql, status, id); err != nil {
		return err
	}

	return nil
}
