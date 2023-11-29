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

type Template struct {
	templateService service.Template
	log             logging.Logger
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

		r.Route(templatesID, func(r chi.Router) {
			r.Patch("/", t.Update)
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

	var template *models.Template
	if err = json.Unmarshal(b, &template); err != nil {
		t.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: already exists
	ctx := context.Background()
	if err = t.templateService.Create(ctx, template); err != nil {
		t.log.Error("cannot create template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t Template) Update(w http.ResponseWriter, r *http.Request) {}

func (t Template) Delete(w http.ResponseWriter, r *http.Request) {}
