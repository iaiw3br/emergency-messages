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

//go:generate mockgen --destination=./mock_store/template.go template Templater
type Templater interface {
	Create(ctx context.Context, template *models.Template) error
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.Template, error)
}

func NewTemplate(db *pgx.Conn, log logging.Logger) Templater {
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
	return err
}

func (t Template) Update(ctx context.Context, template *models.Template) error {
	sql := `
		UPDATE templates 
		SET subject = $1, 
		    text = $2;`
	_, err := t.db.Exec(ctx, sql, template.Subject, template.Text)
	// TODO: return not found error,
	// if errors.Is(err, pgx.ErrNoRows) {
	//
	// }
	return err
}

func (t Template) Delete(ctx context.Context, id uint64) error {
	sql := `DELETE FROM templates WHERE id = $1`
	_, err := t.db.Exec(ctx, sql, id)
	// TODO: return not found error
	// if errors.Is(err, pgx.ErrNoRows) {
	//
	// }
	return err
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
	// TODO: return not found error
	// if errors.Is(err, pgx.ErrNoRows) {
	//
	// }
	if err != nil {
		return nil, err
	}
	return template, nil

}
