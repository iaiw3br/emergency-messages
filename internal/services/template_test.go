package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTemplate_Create(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mock_services.NewMockTemplateStore(controller)
	ctx := context.Background()
	service := NewTemplate(store, log)

	t.Run("when all data have then no error", func(t *testing.T) {
		template := &models.TemplateCreate{
			Subject: "1",
			Text:    "2",
		}

		templateStore := &models.TemplateEntity{
			Subject: "1",
			Text:    "2",
		}

		store.
			EXPECT().
			Create(ctx, templateStore).
			Return(nil)

		err := service.Create(ctx, template)
		assert.NoError(t, err)
	})
	t.Run("when subject is empty then error", func(t *testing.T) {
		template := &models.TemplateCreate{
			Text: "2",
		}

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
	t.Run("when text is empty then error", func(t *testing.T) {
		template := &models.TemplateCreate{
			Subject: "2",
		}

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
	t.Run("when error in update then error", func(t *testing.T) {
		template := &models.TemplateCreate{
			Text:    "2",
			Subject: "2",
		}
		templateStore := &models.TemplateEntity{
			Subject: template.Subject,
			Text:    template.Text,
		}

		store.
			EXPECT().
			Create(ctx, templateStore).
			Return(errors.New("")).
			AnyTimes()

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
}

func TestNewTemplate(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mock_services.NewMockTemplateStore(controller)
	res := NewTemplate(store, log)
	assert.NotNil(t, res)
	assert.Equal(t, store, res.templateStore)
	assert.Equal(t, log, res.log)
}

func TestTemplate_Delete(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	store := mock_services.NewMockTemplateStore(controller)

	t.Run("when stores doesn't return error then no error", func(t *testing.T) {
		uidStr := "9dfc0a1d-7582-40eb-bc50-53a973bd1dbf"
		uid, err := uuid.Parse(uidStr)
		assert.NoError(t, err)

		store.
			EXPECT().
			Delete(ctx, uid, gomock.Any()).
			Return(nil)

		service := NewTemplate(store, log)

		err = service.Delete(ctx, uidStr)
		assert.NoError(t, err)
	})
	t.Run("when stores returns error then error", func(t *testing.T) {
		uidStr := "9dfc0a1d-7582-40eb-bc50-53a973bd1dbf"
		uid, err := uuid.Parse(uidStr)
		assert.NoError(t, err)

		store.
			EXPECT().
			Delete(ctx, uid, gomock.Any()).
			Return(errors.New(""))

		service := NewTemplate(store, log)

		err = service.Delete(ctx, uidStr)
		assert.Error(t, err)
	})
}

func TestTemplate_Update(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	store := mock_services.NewMockTemplateStore(controller)

	t.Run("when all data have then no error", func(t *testing.T) {
		uidStr := "9dfc0a1d-7582-40eb-bc50-53a973bd1dbf"
		uid, err := uuid.Parse(uidStr)
		assert.NoError(t, err)
		updateTemplate := &models.TemplateUpdate{
			ID:      uidStr,
			Subject: "1",
			Text:    "2",
		}
		templateStore := &models.TemplateEntity{
			ID:      uid,
			Subject: updateTemplate.Subject,
			Text:    updateTemplate.Text,
		}

		store.
			EXPECT().
			Update(ctx, templateStore).
			Return(nil)

		service := NewTemplate(store, log)

		err = service.Update(ctx, updateTemplate)
		assert.NoError(t, err)
	})
	t.Run("when stores has error then error", func(t *testing.T) {
		uidStr := "9dfc0a1d-7582-40eb-bc50-53a973bd1dbf"
		uid, err := uuid.Parse(uidStr)
		assert.NoError(t, err)
		updateTemplate := &models.TemplateUpdate{
			ID:      uidStr,
			Subject: "1",
			Text:    "2",
		}
		templateStore := &models.TemplateEntity{
			ID:      uid,
			Subject: updateTemplate.Subject,
			Text:    updateTemplate.Text,
		}

		store.
			EXPECT().
			Update(ctx, templateStore).
			Return(errors.New(""))

		service := NewTemplate(store, log)

		err = service.Update(ctx, updateTemplate)
		assert.Error(t, err)
	})
	t.Run("when stores has error then error", func(t *testing.T) {
		updateTemplate := &models.TemplateUpdate{
			Subject: "1",
			Text:    "2",
		}

		service := NewTemplate(store, log)

		err := service.Update(ctx, updateTemplate)
		assert.Error(t, err)
	})
}
