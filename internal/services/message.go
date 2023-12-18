package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	"projects/emergency-messages/internal/providers/sms/twil"
	"runtime"
	"sync"
)

type MessageService struct {
	messageStore  Message
	templateStore TemplateStore
	userStore     User
	log           logging.Logger
	sender        providers.Sender
}

type MessageStatus string

const (
	Created   MessageStatus = "created"
	Delivered MessageStatus = "delivered"
)

type Message interface {
	Create(ctx context.Context, m *models.MessageEntity) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error
}

func NewMessage(messageStore Message, templateStore TemplateStore, userStore User, log logging.Logger) *MessageService {
	return &MessageService{
		messageStore:  messageStore,
		templateStore: templateStore,
		userStore:     userStore,
		log:           log,
	}
}

func (s *MessageService) Send(ctx context.Context, message models.CreateMessage) error {
	template, err := s.templateStore.GetByID(ctx, message.TemplateID)
	if err != nil {
		s.log.Error(err)
		return err
	}

	usersStore, err := s.userStore.FindByCity(ctx, message.City)
	if err != nil {
		s.log.Error(err)
		return err
	}

	text := fmt.Sprintf(template.Text, message.City, message.Strength)
	newMessage := models.Message{
		Subject: template.Subject,
		Text:    text,
		Status:  models.Created,
	}

	usersCh := make(chan *models.User)
	var wg sync.WaitGroup
	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go s.send(ctx, usersCh, newMessage, &wg)
	}

	users, err := s.transformUsersStoreToUsers(usersStore)
	if err != nil {
		s.log.Error(err)
		return err
	}

	go sendUsersToUsersChannel(users, usersCh)
	wg.Wait()

	return nil
}

func (s *MessageService) send(ctx context.Context, usersCh <-chan *models.User, newMessage models.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range usersCh {
		newMessage.UserID = user.ID

		storeModel, err := s.transformMessageToStoreModel(newMessage)
		if err != nil {
			s.log.Error(err)
			continue
		}
		if err = s.messageStore.Create(ctx, storeModel); err != nil {
			s.log.Error(err)
			continue
		}

		// sms
		if user.MobilePhone != "" {
			sms := twil.NewMobileTwilClient(s.log)
			if err = sms.Send(newMessage, user.MobilePhone); err != nil {
				s.log.Error(err)
				continue
			}
		}

		// email
		if user.Email != "" {
			email := mail_gun.NewEmailMailgClient(s.log)
			if err = email.Send(newMessage, user.Email); err != nil {
				s.log.Error(err)
				continue
			}
		}

		if err = s.messageStore.UpdateStatus(ctx, storeModel.ID, models.Delivered); err != nil {
			s.log.Error(err)
			continue
		}
	}
}

func (s *MessageService) transformMessageToStoreModel(m models.Message) (*models.MessageEntity, error) {
	storeModel := &models.MessageEntity{
		Subject: m.Subject,
		Text:    m.Text,
		Status:  m.Status,
		UserID:  m.UserID,
	}
	return storeModel, nil
}

func (s *MessageService) transformUsersStoreToUsers(usersStore []models.UserEntity) ([]*models.User, error) {
	users := make([]*models.User, 0, len(usersStore))
	for _, u := range usersStore {
		user := &models.User{
			ID:          u.ID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			MobilePhone: u.MobilePhone,
			Email:       u.Email,
			City:        u.City,
		}
		users = append(users, user)
	}
	return users, nil
}

func sendUsersToUsersChannel(users []*models.User, usersCh chan<- *models.User) {
	for _, u := range users {
		usersCh <- u
	}
	close(usersCh)
}
