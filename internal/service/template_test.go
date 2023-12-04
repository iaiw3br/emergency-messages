package service

import (
	"context"
	"errors"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	mock_store "github.com/emergency-messages/internal/store/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTemplate_Create(t *testing.T) {
	t.Run("when all data have then no error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplater(controller)
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
			Return(uint64(1), nil)

		id, err := service.Create(ctx, template)
		assert.NoError(t, err)
		assert.NotNil(t, id)
	})
	t.Run("when subject is empty then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplater(controller)
		template := &models.TemplateCreate{
			Text: "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		id, err := service.Create(ctx, template)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
	})
	t.Run("when text is empty then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplater(controller)
		template := &models.TemplateCreate{
			Subject: "2",
		}
		ctx := context.Background()
		log := logging.New()
		service := NewTemplate(store, log)

		id, err := service.Create(ctx, template)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
	})
	t.Run("when error in update then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		store := mock_store.NewMockTemplater(controller)
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
			Return(uint64(0), errors.New("")).
			AnyTimes()

		id, err := service.Create(ctx, template)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
	})
}

func TestNewTemplate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logging.New()
	store := mock_store.NewMockTemplater(controller)
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
		store := mock_store.NewMockTemplater(controller)

		store.
			EXPECT().
			Delete(ctx, uint64(1)).
			Return(nil)

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Delete(ctx, uint64(1))
		assert.NoError(t, err)
	})
	t.Run("when store returns error then error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplater(controller)

		store.
			EXPECT().
			Delete(ctx, uint64(1)).
			Return(errors.New(""))

		log := logging.New()
		service := NewTemplate(store, log)

		err := service.Delete(ctx, uint64(1))
		assert.Error(t, err)
	})
}

func TestTemplate_Update(t *testing.T) {
	t.Run("when all data have then no error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		ctx := context.Background()
		store := mock_store.NewMockTemplater(controller)
		updateTemplate := &models.TemplateUpdate{
			ID:      1,
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
		store := mock_store.NewMockTemplater(controller)
		updateTemplate := &models.TemplateUpdate{
			ID:      1,
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
		store := mock_store.NewMockTemplater(controller)
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
