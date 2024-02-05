package controllers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
)

type TemplateService interface {
	Create(ctx context.Context, template *models.TemplateCreate) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, template *models.TemplateUpdate) error
}

type Template struct {
	templateService TemplateService
	log             logging.Logger
}

type templateUpdate struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type templateCreate struct {
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func NewTemplate(templateService TemplateService, log logging.Logger) *Template {
	return &Template{
		templateService: templateService,
		log:             log,
	}
}

func (t Template) Create(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.log.Error("cannot read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var temp *templateCreate
	if err = json.Unmarshal(b, &temp); err != nil {
		t.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newTemplate := &models.TemplateCreate{
		Subject: temp.Subject,
		Text:    temp.Text,
	}

	ctx := context.Background()
	err = t.templateService.Create(ctx, newTemplate)
	if assertError(err, w) {
		t.log.Error("Template.Create() error:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t Template) Update(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.log.Error("cannot read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var template *templateUpdate
	if err = json.Unmarshal(b, &template); err != nil {
		t.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := &models.TemplateUpdate{
		ID:      template.ID,
		Subject: template.Subject,
		Text:    template.Text,
	}

	ctx := context.Background()
	err = t.templateService.Update(ctx, u)
	if assertError(err, w) {
		t.log.Error("Template.Update() error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t Template) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := chi.URLParam(r, "id")

	err := t.templateService.Delete(ctx, id)
	if assertError(err, w) {
		t.log.Error("Template.Delete() error:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
