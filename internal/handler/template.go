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
	"strconv"
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

	var template *models.TemplateCreate
	if err = json.Unmarshal(b, &template); err != nil {
		t.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	newTemplated, err := t.templateService.Create(ctx, template)
	if err != nil {
		t.log.Error("cannot create template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	templateBytes, err := json.Marshal(newTemplated)
	if err = json.Unmarshal(b, &template); err != nil {
		t.log.Error("cannot marshal template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(templateBytes)
	w.WriteHeader(http.StatusCreated)
}

func (t Template) Update(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.log.Error("cannot read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var template *models.TemplateUpdate
	if err = json.Unmarshal(b, &template); err != nil {
		t.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := t.templateService.Update(ctx, template); err != nil {
		t.log.Error("cannot create template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t Template) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		t.log.Errorf("cannot transform id:%s to int", idStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = t.templateService.Delete(ctx, uint64(id)); err != nil {
		t.log.Error("template service delete return error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
