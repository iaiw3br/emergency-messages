package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

type MessageService struct {
	messageStore  Message
	templateStore TemplateStore
	userStore     User
	log           *slog.Logger
	sender        *providers.SendManager
}

type Message interface {
	Create(ctx context.Context, m *models.MessageEntity) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus) error
}

func NewMessage(messageStore Message, templateStore TemplateStore, userStore User, sender *providers.SendManager, log *slog.Logger) *MessageService {
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
		s.log.With(slog.Any("templateID", message.TemplateID)).
			Error("getting template", err)
		if err == sql.ErrNoRows {
			return errorx.ErrNotFound
		}
		return errorx.ErrInternal
	}

	usersStore, err := s.userStore.FindByCity(ctx, message.City)
	if err != nil {
		s.log.With(slog.Any("city", message.City)).
			Error("finding user by city", err)
		if err == sql.ErrNoRows {
			return errorx.ErrNotFound
		}
		return errorx.ErrInternal
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
		s.log.Error("transforming store model to user", err)
		return err
	}

	go sendUsersToUsersChannel(users, usersCh)
	wg.Wait()

	return nil
}

func (s *MessageService) send(ctx context.Context, usersCh <-chan *models.User, newMessage models.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range usersCh {
		for _, contact := range user.Contacts {
			if !contact.IsActive {
				continue
			}
			newMessage.UserID = user.ID

			storeModel, err := s.transformMessageToStoreModel(newMessage)
			if err != nil {
				s.log.Error("transforming message to store model", err)
				continue
			}

			if err = s.messageStore.Create(ctx, storeModel); err != nil {
				s.log.Error("creating message", err)
				continue
			}

			if err = s.sender.Send(newMessage, contact); err != nil {
				s.log.Error("sending message", err)
				continue
			}

			if err = s.messageStore.UpdateStatus(ctx, storeModel.ID, models.Delivered); err != nil {
				s.log.Error("updating message", err)
				continue
			}
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
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Contacts:  u.Contacts,
			City:      u.City,
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
