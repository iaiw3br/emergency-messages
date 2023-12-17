package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"

	"github.com/go-chi/chi/v5"
)

const messages = "/messages"

type MessageController struct {
	messageService *services.MessageService
	log            logging.Logger
}

func NewMessage(messageService *services.MessageService, log logging.Logger) *MessageController {
	return &MessageController{
		messageService: messageService,
		log:            log,
	}
}

func (m MessageController) Register(r *chi.Mux) {
	r.Post(messages, m.Send)
}

func (m MessageController) Send(w http.ResponseWriter, r *http.Request) {
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
	if err := m.messageService.Send(ctx, message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
