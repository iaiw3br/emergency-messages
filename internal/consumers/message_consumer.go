// Package consumers
// Handles reading messages from Kafka queue
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"projects/emergency-messages/internal/models"
)

type Consumer struct {
	sender       MessageSender
	messageStore MessageUpdater
	log          *slog.Logger
}

type MessageSender interface {
	Send(message models.MessageConsumer) error
}

type MessageUpdater interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error
}

func New(sender MessageSender, messageStore MessageUpdater, log *slog.Logger) *Consumer {
	return &Consumer{
		sender:       sender,
		messageStore: messageStore,
		log:          log,
	}
}

// Send sends the message
func (c *Consumer) Send(messageBytes []byte) {
	var message models.MessageConsumer
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		c.log.With(slog.Any("message", string(messageBytes))).
			Error("unmarshalling message", err)
		return
	}

	if err := c.sender.Send(message); err != nil {
		c.log.With(slog.Any("message", message)).
			Error("sending message", err)
	}
}

// UpdateMessageStatus updates the status of the message
func (c *Consumer) UpdateMessageStatus(messageBytesID []byte, status models.MessageStatus) {
	fmt.Println("UpdateMessageStatus", string(messageBytesID), status)

	messageID, err := uuid.Parse(string(messageBytesID))
	if err != nil {
		c.log.With(slog.Any("message id", string(messageBytesID))).
			Error("parsing uuid", err)
		return
	}

	ctx := context.Background()
	if err := c.messageStore.UpdateStatus(ctx, messageID, status); err != nil {
		c.log.With(
			slog.Any("message id", messageID),
			slog.Any("status", status),
		).Error("updating message status", err)
		return
	}
}
