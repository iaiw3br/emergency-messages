package service

import (
	"context"
	"errors"
	mock_store "projects/emergency-messages/internal/store/postgres/mock"
	"testing"

	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTemplate_Create(t *testing.T) {
	t.Run("when all data have then no error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplateStore(controller)
		template := &models.TemplateCreate{
			Subject: "1",
			Text:    "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		store.
			EXPECT().
			Create(ctx, template).
			Return(nil)

		err := service.Create(ctx, template)
		assert.NoError(t, err)
	})
	t.Run("when subject is empty then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplateStore(controller)
		template := &models.TemplateCreate{
			Text: "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
	t.Run("when text is empty then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplateStore(controller)
		template := &models.TemplateCreate{
			Subject: "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
	t.Run("when error in update then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplateStore(controller)
		template := &models.TemplateCreate{
			Text:    "2",
			Subject: "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		store.
			EXPECT().
			Create(ctx, template).
			Return(errors.New("")).
			AnyTimes()

		err := service.Create(ctx, template)
		assert.Error(t, err)
	})
}

func TestNewTemplate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logging.New()
	store := mock_store.NewMockTemplateStore(controller)
	res := NewTemplate(store, log)
	assert.NotNil(t, res)
	assert.Equal(t, store, res.templateStore)
	assert.Equal(t, log, res.log)
}

func TestTemplate_Delete(t *testing.T) {
	t.Run("when store doesn't return error then no error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplateStore(controller)

		store.
			EXPECT().
			Delete(ctx, "11111").
			Return(nil)

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Delete(ctx, "11111")
		assert.NoError(t, err)
	})
	t.Run("when store returns error then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplateStore(controller)

		store.
			EXPECT().
			Delete(ctx, "11111").
			Return(errors.New(""))

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Delete(ctx, "11111")
		assert.Error(t, err)
	})
}

func TestTemplate_Update(t *testing.T) {
	t.Run("when all data have then no error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplateStore(controller)
		updateTemplate := &models.TemplateUpdate{
			ID:      "1",
			Subject: "1",
			Text:    "2",
		}

		store.
			EXPECT().
			Update(ctx, updateTemplate).
			Return(nil)

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Update(ctx, updateTemplate)
		assert.NoError(t, err)
	})
	t.Run("when store has error then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplateStore(controller)
		updateTemplate := &models.TemplateUpdate{
			ID:      "1",
			Subject: "1",
			Text:    "2",
		}

		store.
			EXPECT().
			Update(ctx, updateTemplate).
			Return(errors.New(""))

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Update(ctx, updateTemplate)
		assert.Error(t, err)
	})
	t.Run("when store has error then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplateStore(controller)
		updateTemplate := &models.TemplateUpdate{
			Subject: "1",
			Text:    "2",
		}

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Update(ctx, updateTemplate)
		assert.Error(t, err)
	})
}
