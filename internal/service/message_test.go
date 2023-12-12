package service

import (
	"github.com/emergency-messages/internal/logging"
	mailg "github.com/emergency-messages/internal/providers/email/mailgun"
	mock_store "github.com/emergency-messages/internal/store/postgres/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestNewMessage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	messageStore := mock_store.NewMockMessager(controller)
	templateStore := mock_store.NewMockTemplateStore(controller)
	userStore := mock_store.NewMockUser(controller)

	log := logging.New()
	email := mailg.New(log)

	res := NewMessage(messageStore, templateStore, userStore, email, log)
	assert.NotNil(t, res)
	assert.Equal(t, messageStore, res.messageStore)
	assert.Equal(t, templateStore, res.templateStore)
	assert.Equal(t, userStore, res.userStore)
	assert.Equal(t, log, res.log)
}
