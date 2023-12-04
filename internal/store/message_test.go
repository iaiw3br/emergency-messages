package store

import (
	"context"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDatabase(t)

	// create template
	template := &models.Template{
		Subject: "test subject",
		Text:    "test text",
	}

	templateStore := NewTemplate(db, logging.New())
	id, err := templateStore.Create(ctx, template)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	gotTemplate, err := templateStore.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, gotTemplate)

	// create user
	user := models.User{
		FirstName:   "Mark",
		LastName:    "Smith",
		MobilePhone: "+7843286473",
		Email:       "mark@gmail.com",
		City:        "Perth",
	}
	userStore := NewUserStore(db)
	userCreated, err := userStore.Create(ctx, user)
	assert.NoError(t, err)

	// create message
	newMessage := models.Message{
		Subject: gotTemplate.Subject,
		Text:    gotTemplate.Text,
		Status:  models.Created,
		UserID:  userCreated.ID,
	}

	messageStore := NewMessage(db, logging.New())

	id, err = messageStore.Create(ctx, newMessage)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	gotMessage, err := messageStore.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Subject, gotMessage.Subject)
	assert.Equal(t, newMessage.Text, gotMessage.Text)
	assert.Equal(t, models.Created, gotMessage.Status)
	assert.Equal(t, newMessage.UserID, gotMessage.UserID)
}

func TestMessage_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	db := setupTestDatabase(t)

	// create template
	template := &models.Template{
		Subject: "test subject",
		Text:    "test text",
	}
	templateStore := NewTemplate(db, logging.New())
	id, err := templateStore.Create(ctx, template)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	gotTemplate, err := templateStore.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, gotTemplate)

	// create user
	user := models.User{
		FirstName:   "Mark",
		LastName:    "Smith",
		MobilePhone: "+7843286473",
		Email:       "mark@gmail.com",
		City:        "Perth",
	}
	userStore := NewUserStore(db)
	userCreated, err := userStore.Create(ctx, user)
	assert.NoError(t, err)

	// create message
	newMessage := models.Message{
		Subject: gotTemplate.Subject,
		Text:    gotTemplate.Text,
		Status:  models.Created,
		UserID:  userCreated.ID,
	}

	messageStore := NewMessage(db, logging.New())

	id, err = messageStore.Create(ctx, newMessage)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	err = messageStore.UpdateStatus(ctx, id, models.Delivered)
	assert.NoError(t, err)

	gotMessage, err := messageStore.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Subject, gotMessage.Subject)
	assert.Equal(t, newMessage.Text, gotMessage.Text)
	assert.Equal(t, models.Delivered, gotMessage.Status)
	assert.Equal(t, newMessage.UserID, gotMessage.UserID)
}
