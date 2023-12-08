package service

import (
	"context"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/store"
)

type Template struct {
	templateStore store.Templater
	log           logging.Logger
}

func NewTemplate(templateStore store.Templater, log logging.Logger) Template {
	return Template{
		templateStore: templateStore,
		log:           log,
	}
}

// Create a new template
func (t Template) Create(ctx context.Context, template *models.TemplateCreate) (uint64, error) {
	if err := template.Validate(); err != nil {
		t.log.Error(err)
		return 0, err
	}

	id, err := t.templateStore.Create(ctx, template)
	if err != nil {
		t.log.Errorf("cannot create template %v", template)
		return 0, err
	}

	return id, nil
}

func (t Template) Delete(ctx context.Context, id uint64) error {
	return t.templateStore.Delete(ctx, id)
}

func (t Template) Update(ctx context.Context, template *models.TemplateUpdate) error {
	if err := template.Validate(); err != nil {
		t.log.Error(err)
		return err
	}
	return t.templateStore.Update(ctx, template)
}
