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
func (t Template) Create(ctx context.Context, template *models.Template) error {
	if template.Subject == "" {
		return errors.New("subject is empty")
	}
	if template.Text == "" {
		return errors.New("text is empty")
	}

	if err := t.templateStore.Create(ctx, template); err != nil {
		t.log.Errorf("cannot create template %v", template)
		return err
	}

	return nil
}

func (t Template) Delete() {}

func (t Template) Update() {}
