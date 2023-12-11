package service

import (
	"context"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/providers"
	"runtime"
	"sync"
)

type MessageService struct {
	messageStore  Messager
	templateStore TemplateStore
	userStore     User
	log           logging.Logger
	sender        providers.Sender
}

type Messager interface {
	Create(ctx context.Context, message *models.Message) error
	UpdateStatus(ctx context.Context, id string, status models.MessageStatus) error
}

func NewMessage(messageStore Messager, templateStore TemplateStore, userStore User, sender providers.Sender, log logging.Logger) *MessageService {
	return &MessageService{
		messageStore:  messageStore,
		templateStore: templateStore,
		userStore:     userStore,
		log:           log,
		sender:        sender,
	}
}

func (m MessageService) Send(ctx context.Context, message models.CreateMessage) error {
	template, err := m.templateStore.GetByID(ctx, message.TemplateID)
	if err != nil {
		m.log.Error(err)
		return err
	}

	users, err := m.userStore.FindByCity(ctx, message.City)
	if err != nil {
		m.log.Error(err)
		return err
	}

	text := fmt.Sprintf(template.Text, message.City, message.Strength)
	newMessage := models.Message{
		Subject: template.Subject,
		Text:    text,
		Status:  models.Created,
	}

	usersCh := make(chan models.User)
	var wg sync.WaitGroup
	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go m.send(ctx, usersCh, newMessage, &wg)
	}

	go sendUsersToUsersChannel(users, usersCh)
	wg.Wait()

	return nil
}

func (m MessageService) send(ctx context.Context, usersCh <-chan models.User, newMessage models.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range usersCh {
		newMessage.UserID = user.ID

		if err := m.messageStore.Create(ctx, &newMessage); err != nil {
			m.log.Errorf("cannot create new message: %v", newMessage)
			continue
		}

		if err := m.sender.Send(newMessage, user.Email); err != nil {
			m.log.Errorf("cannot send email to: %s", user.Email)
			continue
		}

		if err := m.messageStore.UpdateStatus(ctx, newMessage.ID, models.Delivered); err != nil {
			m.log.Errorf("cannot update message id: %s to status %s", newMessage.ID, newMessage.Status)
			continue
		}
	}
}

func sendUsersToUsersChannel(users []models.User, usersCh chan<- models.User) {
	for _, u := range users {
		usersCh <- u
	}
	close(usersCh)
}
