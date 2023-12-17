package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"projects/emergency-messages/internal/service"
)

type templateStore struct {
	db *bun.DB
}

func NewTemplate(db *bun.DB) service.TemplateStore {
	return &templateStore{
		db: db,
	}
}

// Create creates the struct of a template in the database.
// It takes in a context, the new struct of the template.
// It returns an error if the create operation fails.
func (s *templateStore) Create(ctx context.Context, t *service.TemplateEntity) error {
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
func (s *templateStore) Update(ctx context.Context, t *service.TemplateEntity) error {
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
		Model(service.TemplateEntity{}).
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
		return fmt.Errorf("deleting template: couldn't find template with id: %s", id)
	}
	return nil
}

// GetByID retrieves a template from the database by its ID.
// It takes in a context and the ID of the template.
// It returns the template and an error if the retrieval operation fails.
func (s *templateStore) GetByID(ctx context.Context, id uuid.UUID) (*service.TemplateEntity, error) {
	entity := &service.TemplateEntity{ID: id}
	err := s.db.
		NewSelect().
		Model(entity).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting by id template: couldn't get template with id: %s. Error: %w", id, err)
	}

	template := &service.TemplateEntity{
		ID:      entity.ID,
		Subject: entity.Subject,
		Text:    entity.Text,
	}
	return template, nil
}
