package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"

	"github.com/go-chi/chi/v5"
)

const (
	templates   = "/templates"
	templatesID = "/{id}"
)

type TemplateService interface {
	Create(ctx context.Context, template *models.TemplateCreate) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, template *models.TemplateUpdate) error
}

type Template struct {
	templateService services.TemplateService
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

func NewTemplate(templateService services.TemplateService, log logging.Logger) Template {
	return Template{
		templateService: templateService,
		log:             log,
	}
}

func (t Template) Register(r *chi.Mux) {
	r.Route(templates, func(r chi.Router) {
		r.Post("/", t.Create)
		r.Patch("/", t.Update)

		r.Route(templatesID, func(r chi.Router) {
			r.Delete("/", t.Delete)
		})
	})
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
	if httpCode := assertError(err); httpCode != 0 {
		t.log.Error("Template.Create() error:", err)
		http.Error(w, http.StatusText(httpCode), httpCode)
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
	if httpCode := assertError(err); httpCode != 0 {
		t.log.Error("Template.Update() error:", err)
		http.Error(w, http.StatusText(httpCode), httpCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t Template) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := chi.URLParam(r, "id")

	err := t.templateService.Delete(ctx, id)
	if httpCode := assertError(err); httpCode != 0 {
		t.log.Error("Template.Delete() error:", err)
		http.Error(w, http.StatusText(httpCode), httpCode)
		return
	}
	w.WriteHeader(http.StatusOK)
}
