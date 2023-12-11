package handler

import (
	"context"
	"encoding/json"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const (
	templates   = "/templates"
	templatesID = "/:id"
)

type TemplateService interface {
	Create(ctx context.Context, template *models.TemplateCreate) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, template *models.TemplateUpdate) error
}

type Template struct {
	templateService service.Template
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

func NewTemplate(templateService service.Template, log logging.Logger) Template {
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
	if err = t.templateService.Create(ctx, newTemplate); err != nil {
		t.log.Error("cannot create template")
		w.WriteHeader(http.StatusBadRequest)
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
	if err := t.templateService.Update(ctx, u); err != nil {
		t.log.Error("cannot create template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t Template) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.URL.Query().Get("id")

	if err := t.templateService.Delete(ctx, id); err != nil {
		t.log.Error("template service delete return error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
