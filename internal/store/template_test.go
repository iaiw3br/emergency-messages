package store

import (
	"context"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/pkg/tests"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestDatabase(t *testing.T) *pgx.Conn {
	t.Helper()
	db, err := tests.SetupTestDatabase()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	return db
}

func TestTemplate_Create(t *testing.T) {
	db := setupTestDatabase(t)
	ctx := context.Background()

	wantTemplate := &models.TemplateCreate{
		Subject: "MSCH",
		Text:    "be careful",
	}
	store := NewTemplate(db, logging.New())
	id, err := store.Create(ctx, wantTemplate)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	gotTemplate, err := store.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, wantTemplate.Subject, gotTemplate.Subject)
	assert.Equal(t, wantTemplate.Text, gotTemplate.Text)
}

func TestTemplate_Update(t *testing.T) {
	t.Run("when create, update template then no error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		db := setupTestDatabase(t)

		template := &models.TemplateCreate{
			Subject: "MSCH",
			Text:    "be careful",
		}
		store := NewTemplate(db, logging.New())
		id, err := store.Create(ctx, template)
		assert.NoError(t, err)
		assert.NotNil(t, id)

		templateToUpdate := &models.TemplateUpdate{
			ID:      id,
			Subject: "new subject",
			Text:    "new text",
		}
		// act
		err = store.Update(ctx, templateToUpdate)
		assert.NoError(t, err)

		// assert
		gotTemplate, err := store.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, templateToUpdate.Subject, gotTemplate.Subject)
		assert.Equal(t, templateToUpdate.Text, gotTemplate.Text)
	})
	t.Run("when update template which doesn't exist in database then error", func(t *testing.T) {
		ctx := context.Background()
		db := setupTestDatabase(t)
		store := NewTemplate(db, logging.New())

		templateToUpdate := &models.TemplateUpdate{
			Subject: "new subject",
			Text:    "new text",
		}
		err := store.Update(ctx, templateToUpdate)
		assert.Error(t, err)
	})
}

func TestTemplate_Delete(t *testing.T) {
	t.Run("when delete template then no error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		db := setupTestDatabase(t)

		template := &models.TemplateCreate{
			Subject: "MSCH",
			Text:    "be careful",
		}
		store := NewTemplate(db, logging.New())
		id, err := store.Create(ctx, template)
		assert.NoError(t, err)
		assert.NotNil(t, id)

		// act
		err = store.Delete(ctx, id)
		assert.NoError(t, err)
	})
	t.Run("when template doesn't exist in database then error", func(t *testing.T) {
		ctx := context.Background()
		db := setupTestDatabase(t)
		store := NewTemplate(db, logging.New())

		err := store.Delete(ctx, 999999999)
		assert.Error(t, err)
	})
}
