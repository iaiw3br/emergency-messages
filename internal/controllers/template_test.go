package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	mock_controllers "projects/emergency-messages/internal/controllers/mocks"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"testing"
)

func TestTemplate_Create(t *testing.T) {
	t.Run("when data is valid then no error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateCreate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateCreate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(nil)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Create)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("when body is invalid then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		service := mock_controllers.NewMockTemplateService(controller)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Create)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := `{"subject":"valid", "text":"valid"`
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("when service return validation error then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateCreate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateCreate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(errorx.ErrValidation)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Create)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("when service return internal error then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateCreate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateCreate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(errorx.ErrInternal)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Create)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestTemplate_Update(t *testing.T) {
	t.Run("when data is valid then no error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateUpdate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateUpdate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(nil)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Update)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("when body is invalid then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		service := mock_controllers.NewMockTemplateService(controller)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Update)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := `{"subject":"valid", "text":"valid"`
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("when service return validation error then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateUpdate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateUpdate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(errorx.ErrValidation)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Update)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("when service return internal error then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		// Prepare the request body
		body := &templateUpdate{
			Subject: "SomeValidSubject",
			Text:    "SomeValidText",
		}

		newTemplate := &models.TemplateUpdate{
			Subject: body.Subject,
			Text:    body.Text,
		}
		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(errorx.ErrInternal)

		log := logging.New()
		tmpl := NewTemplate(service, log)

		r.Post("/templates", tmpl.Update)

		ts := httptest.NewServer(r)
		defer ts.Close()

		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/templates", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
