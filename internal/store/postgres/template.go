package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/service"
	"github.com/uptrace/bun"
	"time"
)

type templateStore struct {
	db *bun.DB
}

type templateEntity struct {
	bun.BaseModel `bun:"table:templates,alias:t"`
	ID            string    `bun:"type:uuid"`
	Subject       string    `bun:"subject,notnull"`
	Text          string    `bun:"text,notnull"`
	Created       time.Time `bun:"created,notnull"`
}

func NewTemplate(db *bun.DB) service.TemplateStore {
	return &templateStore{
		db: db,
	}
}

// Create template
func (s *templateStore) Create(ctx context.Context, t *models.TemplateCreate) error {
	entity := templateEntity{
		ID:      t.ID,
		Subject: t.Subject,
		Text:    t.Text,
		Created: t.Created,
	}
	_, err := s.db.NewInsert().Model(&entity).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Update template
func (s *templateStore) Update(ctx context.Context, t *models.TemplateUpdate) error {
	entity := templateEntity{
		ID:      t.ID,
		Subject: t.Subject,
		Text:    t.Text,
	}
	exec, err := s.db.NewUpdate().
		Model(&entity).
		Where("id = ?", entity.ID).
		Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New(fmt.Sprintf("template was not found by id: %s", t.ID))
	}

	return nil
}

// Delete template by id
func (s *templateStore) Delete(ctx context.Context, id string) error {
	entity := templateEntity{ID: id}
	exec, err := s.db.NewDelete().
		Model(&entity).
		Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New(fmt.Sprintf("template was not found by id: %s", id))
	}
	return nil
}

// GetByID find by ID and return store
func (s *templateStore) GetByID(ctx context.Context, id string) (*models.Template, error) {
	entity := new(templateEntity)
	err := s.db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	template := &models.Template{
		ID:      entity.ID,
		Subject: entity.Subject,
		Text:    entity.Text,
	}
	return template, nil
}
