package workers

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"projects/emergency-messages/internal/models"
)

type UpdateStatusMessage struct {
	messageStore MessageUpdater
	log          *slog.Logger
}

type MessageUpdater interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error
}

func (u *UpdateStatusMessage) UpdateStatus(messageID uuid.UUID, status models.MessageStatus) error {
	ctx := context.Background()
	return u.messageStore.UpdateStatus(ctx, messageID, status)
}
