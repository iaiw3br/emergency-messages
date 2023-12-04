package store

import (
	"context"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type Template struct {
	db  *pgx.Conn
	log logging.Logger
}

type Templater interface {
	Create(ctx context.Context, template *models.TemplateCreate) (uint64, error)
	Update(ctx context.Context, template *models.TemplateUpdate) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.Template, error)
}

func NewTemplate(db *pgx.Conn, log logging.Logger) Templater {
	return Template{
		db:  db,
		log: log,
	}
}

func (t Template) Create(ctx context.Context, template *models.TemplateCreate) (uint64, error) {
	row := t.db.QueryRow(ctx, `
		INSERT INTO templates (subject, text) 
		VALUES ($1, $2) 
		RETURNING id;
	`, template.Subject, template.Text)
	var result uint64
	if err := row.Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func (t Template) Update(ctx context.Context, template *models.TemplateUpdate) error {
	sql := `
		UPDATE templates 
		SET subject = $1, 
		    text = $2
		WHERE id = $3;`
	tag, err := t.db.Exec(ctx, sql, template.Subject, template.Text, template.ID)
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (t Template) Delete(ctx context.Context, id uint64) error {
	sql := `DELETE FROM templates WHERE id = $1`
	tag, err := t.db.Exec(ctx, sql, id)
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	fmt.Println(tag)
	return err
}

// GetByID find by ID and return Template
func (t Template) GetByID(ctx context.Context, id uint64) (*models.Template, error) {
	sql := `
		SELECT id, subject, text
		FROM templates
		WHERE id = $1;
	`
	template := &models.Template{}
	err := t.db.QueryRow(ctx, sql, id).Scan(&template.ID, &template.Subject, &template.Text)
	if err != nil {
		return nil, err
	}
	return template, nil
}
