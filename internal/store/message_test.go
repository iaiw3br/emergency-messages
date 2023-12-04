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
	templateCreated, err := templateStore.Create(ctx, template)
	assert.NoError(t, err)

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
		Subject: templateCreated.Subject,
		Text:    templateCreated.Text,
		Status:  models.Created,
		UserID:  userCreated.ID,
	}

	messageStore := NewMessage(db, logging.New())

	messageCreated, err := messageStore.Create(ctx, newMessage)
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Subject, messageCreated.Subject)
	assert.Equal(t, newMessage.Text, messageCreated.Text)
	assert.Equal(t, newMessage.Status, messageCreated.Status)
	assert.Equal(t, newMessage.UserID, messageCreated.UserID)
	assert.NotNil(t, messageCreated.ID)
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
	templateCreated, err := templateStore.Create(ctx, template)
	assert.NoError(t, err)

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
		Subject: templateCreated.Subject,
		Text:    templateCreated.Text,
		Status:  models.Created,
		UserID:  userCreated.ID,
	}

	messageStore := NewMessage(db, logging.New())

	messageCreated, err := messageStore.Create(ctx, newMessage)
	assert.NoError(t, err)

	messageCreated.Deliver()
	err = messageStore.UpdateStatus(ctx, messageCreated.ID, messageCreated.Status)
	assert.NoError(t, err)

	gotMessage, err := messageStore.GetByID(ctx, messageCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Subject, gotMessage.Subject)
	assert.Equal(t, newMessage.Text, gotMessage.Text)
	assert.Equal(t, models.Delivered, gotMessage.Status)
	assert.Equal(t, newMessage.UserID, gotMessage.UserID)
	assert.NotNil(t, messageCreated.ID)
}
