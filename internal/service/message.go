package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers"
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

type MessageEntity struct {
	bun.BaseModel `bun:"table:messages,alias:m"`
	ID            uuid.UUID     `bun:"type:uuid,primarykey"`
	Subject       string        `bun:"subject,notnull"`
	Text          string        `bun:"text,notnull"`
	Status        MessageStatus `bun:"status,notnull"`
	UserID        uuid.UUID     `bun:"user_id,notnull"`
}

type MessageStatus string

const (
	Created   MessageStatus = "created"
	Delivered MessageStatus = "delivered"
)

type Message interface {
	Create(ctx context.Context, m *MessageEntity) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error
}

func NewMessage(messageStore Message, templateStore TemplateStore, userStore User, sender providers.Sender, log logging.Logger) *MessageService {
	return &MessageService{
		messageStore:  messageStore,
		templateStore: templateStore,
		userStore:     userStore,
		log:           log,
		sender:        sender,
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

		if err = s.sender.Send(newMessage, user.Email); err != nil {
			s.log.Error(err)
			continue
		}

		if err = s.messageStore.UpdateStatus(ctx, storeModel.ID, models.Delivered); err != nil {
			s.log.Error(err)
			continue
		}
	}
}

func (s *MessageService) transformMessageToStoreModel(m models.Message) (*MessageEntity, error) {
	storeModel := &MessageEntity{
		ID:      m.ID,
		Subject: m.Subject,
		Text:    m.Text,
		Status:  MessageStatus(m.Status),
		UserID:  m.UserID,
	}
	return storeModel, nil
}

func (s *MessageService) transformUsersStoreToUsers(usersStore []UserEntity) ([]*models.User, error) {
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
