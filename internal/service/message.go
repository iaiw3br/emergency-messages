package service

import (
	"context"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/store"
	"github.com/mailgun/mailgun-go"
	"os"
)

type MessageService struct {
	messageStore  store.Message
	templateStore store.Template
	userStore     store.User
	log           logging.Logger
}

func NewMessage(messageStore store.Message, templateStore store.Template, userStore store.User, log logging.Logger) *MessageService {
	return &MessageService{
		messageStore:  messageStore,
		templateStore: templateStore,
		userStore:     userStore,
		log:           log,
	}
}

func (m MessageService) Send(ctx context.Context, message models.CreateMessage) error {
	template, err := m.templateStore.GetByID(ctx, uint64(message.TemplateID))
	if err != nil {
		m.log.Errorf("cannot find template by id: %d", message.TemplateID)
		return err
	}

	users, err := m.userStore.FindByCity(ctx, message.City)
	if err != nil {
		m.log.Errorf("cannot find users by city: %s", message.City)
		return err
	}

	for _, user := range users {
		text := fmt.Sprintf(template.Text, message.City, message.Strength)
		newMessage := models.Message{
			UserID:  user.ID,
			Subject: template.Subject,
			Text:    text,
			Status:  models.Created,
		}

		if err = m.messageStore.Create(ctx, newMessage); err != nil {
			m.log.Errorf("cannot create new message: %v", newMessage)
			return err
		}

		// TODO: add send by sms and add goroutine
		_, err := m.sendByEmail(newMessage, user.Email)
		if err != nil {
			m.log.Errorf("cannot send email to: %s", user.Email)
			return err
		}

		newMessage.Deliver()
		if err = m.messageStore.UpdateStatus(ctx, newMessage.ID, newMessage.Status); err != nil {
			m.log.Errorf("cannot update message: %d to status %s", newMessage.ID, newMessage.Status)
			return err
		}
	}

	return nil
}

func (m MessageService) sendByEmail(newMessage models.Message, email string) (string, error) {
	apiKey := os.Getenv("EMAIL_API_KEY")
	domain := os.Getenv("EMAIL_DOMAIN")
	mg := mailgun.NewMailgun(domain, apiKey)
	msg := mg.NewMessage(
		"Mailgun Sandbox <postmaster@sandboxd2c7ff34f9f943b98ad2246a10bb7d8a.mailgun.org>",
		newMessage.Subject,
		newMessage.Text,
		email,
	)
	_, id, err := mg.Send(msg)
	return id, err
}
