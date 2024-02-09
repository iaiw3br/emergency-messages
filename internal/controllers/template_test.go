package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	mock_controllers "projects/emergency-messages/internal/controllers/mocks"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	"testing"
)

func TestTemplate_Create(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

func TestTemplate_Delete(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("when data is valid then no error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		id := "7d603549-b079-4016-b81e-9e4386c1de21"

		service.EXPECT().
			Delete(ctx, id).
			Return(nil)

		tmpl := NewTemplate(service, log)

		r.Delete(fmt.Sprintf("/templates/%s", id), tmpl.Delete)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/templates/%s", id), nil)
		assert.NoError(t, err)

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("when get not found then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		id := "7d603549-b079-4016-b81e-9e4386c1de21"

		service.EXPECT().
			Delete(ctx, id).
			Return(errorx.ErrNotFound)

		tmpl := NewTemplate(service, log)

		r.Delete(fmt.Sprintf("/templates/%s", id), tmpl.Delete)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/templates/%s", id), nil)
		assert.NoError(t, err)

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("when get internal error then error", func(t *testing.T) {
		r := chi.NewRouter()

		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		service := mock_controllers.NewMockTemplateService(controller)

		id := "7d603549-b079-4016-b81e-9e4386c1de21"

		service.EXPECT().
			Delete(ctx, id).
			Return(errorx.ErrInternal)

		tmpl := NewTemplate(service, log)

		r.Delete(fmt.Sprintf("/templates/%s", id), tmpl.Delete)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/templates/%s", id), nil)
		assert.NoError(t, err)

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
