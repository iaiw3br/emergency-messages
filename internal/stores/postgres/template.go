package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type templateStore struct {
	db *bun.DB
}

func NewTemplate(db *bun.DB) services.TemplateStore {
	return &templateStore{
		db: db,
	}
}

// Create creates the struct of a template in the database.
// It takes in a context, the new struct of the template.
// It returns an error if the create operation fails.
func (s *templateStore) Create(ctx context.Context, t *models.TemplateEntity) error {
	_, err := s.db.
		NewInsert().
		Model(t).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("creating template: couldn't create with: %v. Error: %w", t, err)
	}
	return nil
}

// Update updates the struct of a template in the database.
// It takes in a context, the new struct of the template.
// It returns an error if the update operation fails.
func (s *templateStore) Update(ctx context.Context, t *models.TemplateEntity) error {
	exec, err := s.db.
		NewUpdate().
		Model(t).
		Where("id = ?", t.ID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("updating template: couldn't update with id: %s. Error: %w", t.ID, err)
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return fmt.Errorf("updating template: couldn't get the number of rows affected with id: %s. Error: %w", t.ID, err)
	}
	if affected == 0 {
		return fmt.Errorf("updating template: couldn't find template with id: %s", t.ID)
	}
	return nil
}

// Delete deletes the struct of a template in the database.
// It takes in a context and the ID of the template.
// It returns an error if the delete operation fails.
func (s *templateStore) Delete(ctx context.Context, id uuid.UUID) error {
	exec, err := s.db.
		NewDelete().
		Model(&models.TemplateEntity{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting template: couldn't delete with id: %s. Error: %w", id, err)
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleting template: couldn't get the number of rows affected with id: %s. Error: %w", id, err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetByID retrieves a template from the database by its ID.
// It takes in a context and the ID of the template.
// It returns the template and an error if the retrieval operation fails.
func (s *templateStore) GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateEntity, error) {
	entity := &models.TemplateEntity{ID: id}
	err := s.db.
		NewSelect().
		Model(entity).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting by id template: couldn't get template with id: %s. Error: %w", id, err)
	}

	template := &models.TemplateEntity{
		ID:      entity.ID,
		Subject: entity.Subject,
		Text:    entity.Text,
	}
	return template, nil
}
