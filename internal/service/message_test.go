package service

import (
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	"testing"

	mock_store "projects/emergency-messages/internal/store/postgres/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewMessage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	messageStore := mock_store.NewMockMessager(controller)
	templateStore := mock_store.NewMockTemplateStore(controller)
	userStore := mock_store.NewMockUser(controller)

	log := logging.New()
	email := mail_gun.New(log)

	res := NewMessage(messageStore, templateStore, userStore, email, log)
	assert.NotNil(t, res)
	assert.Equal(t, messageStore, res.messageStore)
	assert.Equal(t, templateStore, res.templateStore)
	assert.Equal(t, userStore, res.userStore)
	assert.Equal(t, log, res.log)
}
