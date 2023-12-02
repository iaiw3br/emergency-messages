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
	return db
}

func TestTemplate_Create(t *testing.T) {
	db := setupTestDatabase(t)

	wantTemplate := &models.Template{
		Subject: "MSCH",
		Text:    "be careful",
	}
	store := NewTemplate(db, logging.New())
	gotTemplate, err := store.Create(context.Background(), wantTemplate)
	assert.NoError(t, err)

	assert.Equal(t, wantTemplate.Subject, gotTemplate.Subject)
	assert.Equal(t, wantTemplate.Text, gotTemplate.Text)
	assert.NotNil(t, gotTemplate.ID)
}

func TestTemplate_Update(t *testing.T) {
	t.Run("when create, update template then no error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		db := setupTestDatabase(t)

		template := &models.Template{
			Subject: "MSCH",
			Text:    "be careful",
		}
		store := NewTemplate(db, logging.New())
		templateCreated, err := store.Create(ctx, template)
		assert.NoError(t, err)

		templateToUpdate := &models.Template{
			ID:      templateCreated.ID,
			Subject: "new subject",
			Text:    "new text",
		}
		// act
		err = store.Update(ctx, templateToUpdate)

		// assert
		assert.NoError(t, err)
		gotTemplate, err := store.GetByID(ctx, templateCreated.ID)
		assert.NoError(t, err)

		assert.Equal(t, templateToUpdate.Subject, gotTemplate.Subject)
		assert.Equal(t, templateToUpdate.Text, gotTemplate.Text)
	})
	t.Run("when update template which doesn't exist in database then error", func(t *testing.T) {
		ctx := context.Background()
		db := setupTestDatabase(t)
		store := NewTemplate(db, logging.New())

		templateToUpdate := &models.Template{
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

		template := &models.Template{
			Subject: "MSCH",
			Text:    "be careful",
		}
		store := NewTemplate(db, logging.New())
		templateCreated, err := store.Create(ctx, template)
		assert.NoError(t, err)

		// act
		err = store.Delete(ctx, templateCreated.ID)
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
