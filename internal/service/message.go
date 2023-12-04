package service

import (
	"context"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	mailg "github.com/emergency-messages/internal/providers/email/mailgun"
	"github.com/emergency-messages/internal/store"
	"runtime"
	"sync"
	"time"
)

type MessageService struct {
	messageStore  store.Messager
	templateStore store.Templater
	userStore     store.User
	log           logging.Logger
	email         *mailg.Client
}

func NewMessage(messageStore store.Messager, templateStore store.Templater, userStore store.User, email *mailg.Client, log logging.Logger) *MessageService {
	return &MessageService{
		messageStore:  messageStore,
		templateStore: templateStore,
		userStore:     userStore,
		log:           log,
		email:         email,
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

		id, err := m.messageStore.Create(ctx, newMessage)
		if err != nil {
			m.log.Errorf("cannot create new message: %v", newMessage)
			continue
		}

		time.Sleep(time.Second * 3)
		// if err := m.email.Send(newMessage, user.Email); err != nil {
		// 	m.log.Errorf("cannot send email to: %s", user.Email)
		// 	continue
		// }

		if err := m.messageStore.UpdateStatus(ctx, id, models.Delivered); err != nil {
			m.log.Errorf("cannot update message: %d to status %s", newMessage.ID, newMessage.Status)
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
