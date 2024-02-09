package services

import (
	"log/slog"
	"os"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	mock_service "projects/emergency-messages/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewMessage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	messageStore := mock_service.NewMockMessage(controller)
	templateStore := mock_service.NewMockTemplateStore(controller)
	userStore := mock_service.NewMockUser(controller)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := providers.New()
	mail := mail_gun.NewEmailMailgClient(log)
	client.AddProvider(mail, models.ContactTypeEmail)

	res := NewMessage(messageStore, templateStore, userStore, client, log)
	assert.NotNil(t, res)
	assert.Equal(t, messageStore, res.messageStore)
	assert.Equal(t, templateStore, res.templateStore)
	assert.Equal(t, userStore, res.userStore)
	assert.Equal(t, log, res.log)
}
