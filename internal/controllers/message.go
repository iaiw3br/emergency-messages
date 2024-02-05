package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"
)

type Message struct {
	messageService *services.MessageService
	log            logging.Logger
}

func NewMessage(messageService *services.MessageService, log logging.Logger) *Message {
	return &Message{
		messageService: messageService,
		log:            log,
	}
}

func (m Message) Send(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		m.log.Error("cannot read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message models.CreateMessage
	if err = json.Unmarshal(b, &message); err != nil {
		m.log.Error("cannot unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = m.messageService.Send(ctx, message)
	if assertError(err, w) {
		m.log.Error("Message.Send() error:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}