package services

import (
	"log/slog"
	"os"
	mock_service "projects/emergency-messages/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewMessage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	templateStore := mock_service.NewMockTemplateStore(controller)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	producer := mock_service.NewMockProducer(controller)

	res := NewMessage(producer, templateStore, log)
	assert.NotNil(t, res)
	assert.Equal(t, templateStore, res.templateStore)
	assert.Equal(t, log, res.log)
}
