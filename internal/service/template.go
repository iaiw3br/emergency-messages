package service

import (
	"context"
	"errors"
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
func (t Template) Create(ctx context.Context, template *models.Template) (uint64, error) {
	if template.Subject == "" {
		err := errors.New("subject is empty")
		t.log.Error(err)
		return 0, err
	}
	if template.Text == "" {
		err := errors.New("text is empty")
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

func (t Template) Delete() {}

func (t Template) Update(ctx context.Context, template *models.Template) error {
	return t.templateStore.Update(ctx, template)
}
