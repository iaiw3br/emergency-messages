package service

import (
	"context"
	"errors"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"time"
)

type Template struct {
	templateStore TemplateStore
	log           logging.Logger
}

type TemplateStore interface {
	Create(ctx context.Context, template *models.TemplateCreate) error
	Update(ctx context.Context, template *models.TemplateUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Template, error)
}

func NewTemplate(templateStore TemplateStore, log logging.Logger) Template {
	return Template{
		templateStore: templateStore,
		log:           log,
	}
}

// Create a new template
func (t Template) Create(ctx context.Context, template *models.TemplateCreate) error {
	if err := template.Validate(); err != nil {
		t.log.Error(err)
		return err
	}

	template.Create(time.Now())

	if err := t.templateStore.Create(ctx, template); err != nil {
		t.log.Errorf("cannot create template %v", template)
		return err
	}

	return nil
}

// Delete template by id
func (t Template) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	return t.templateStore.Delete(ctx, id)
}

// Update template
func (t Template) Update(ctx context.Context, template *models.TemplateUpdate) error {
	if err := template.Validate(); err != nil {
		t.log.Error(err)
		return err
	}
	return t.templateStore.Update(ctx, template)
}
