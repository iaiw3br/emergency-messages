package store

import (
	"context"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type Template struct {
	db  *pgx.Conn
	log logging.Logger
}

func NewTemplate(db *pgx.Conn, log logging.Logger) Template {
	return Template{
		db:  db,
		log: log,
	}
}

func (t Template) Create(ctx context.Context, template *models.Template) error {
	sql := `
		INSERT INTO templates (subject, text) 
		VALUES ($1, $2)
	`
	_, err := t.db.Exec(ctx, sql, template.Subject, template.Text)
	if err != nil {
		return err
	}
	return nil
}

func (t Template) Update() error {
	return nil
}

func (t Template) Delete() error {
	return nil
}

// GetByID find by ID and return Template
func (t Template) GetByID(ctx context.Context, id uint64) (*models.Template, error) {
	sql := `
		SELECT subject, text
		FROM templates
		WHERE id = $1;
	`
	template := &models.Template{}
	err := t.db.QueryRow(ctx, sql, id).Scan(&template.Subject, &template.Text)
	// FIXME: not found,
	if err != nil {
		return nil, err
	}
	return template, nil

}
