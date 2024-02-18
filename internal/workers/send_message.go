package workers

import (
	"context"
	"fmt"
	"log/slog"
	"projects/emergency-messages/internal/models"
)

type Message struct {
	messageStore MessageFinder
	sender       Sender
	producer     Producer
	maxRetries   int
	log          *slog.Logger
}

type MessageFinder interface {
	FindByStatus(ctx context.Context, status models.MessageStatus) ([]models.MessageEntity, error)
}

type Producer interface {
	Delivered(messageID []byte) error
	Failed(messageID []byte) error
}

type Sender interface {
	Send(message models.MessageSend) error
}

func NewSendMessage(messageStore MessageFinder, producer Producer, sender Sender, log *slog.Logger) *Message {
	return &Message{
		messageStore: messageStore,
		producer:     producer,
		sender:       sender,
		log:          log,
	}
}

func (m *Message) Send() {
	fmt.Println("send worker")
	messagesStore, err := m.messageStore.FindByStatus(context.Background(), models.Created)
	if err != nil {
		m.log.Error("finding messages by status", err)
		return
	}

	messages, err := m.transformMessagesStoreToMessages(messagesStore)
	if err != nil {
		m.log.Error("transforming messages store to messages", err)
		return
	}

	// TODO: send message N times with timeout if failed and then update status to failed
	for i, message := range messages {
		b := []byte(message.ID.String())
		if i%2 == 0 {
			if err = m.producer.Failed(b); err != nil {
				m.log.Error("sending failed message", err)
			}
			continue
		}
		// if err = m.sender.Send(message); err != nil {
		// 	if err = m.producer.Failed(message.ID.NodeID()); err != nil {
		// 		m.log.Error("sending failed message", err)
		// 	}
		// 	continue
		// }

		if err = m.producer.Delivered(b); err != nil {
			m.log.Error("sending delivered message", err)
		}
	}
}

func (m *Message) transformMessagesStoreToMessages(messagesStore []models.MessageEntity) ([]models.MessageSend, error) {
	messages := make([]models.MessageSend, 0, len(messagesStore))
	for _, message := range messagesStore {
		newMessage := models.MessageSend{
			ID:      message.ID,
			Subject: message.Subject,
			Text:    message.Text,
			Status:  message.Status,
			UserID:  message.UserID,
			Type:    message.Type,
			Value:   message.Value,
		}
		messages = append(messages, newMessage)
	}
	return messages, nil
}
