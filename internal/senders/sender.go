// Package senders
// Implements logic to send messages to clients
package senders

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	"runtime"
	"sync"
)

type Sender struct {
	messageStore MessageCreator
	userStore    UserFinder
	log          *slog.Logger
}

type MessageCreator interface {
	Create(ctx context.Context, m *models.MessageEntity) error
}

type UserFinder interface {
	FindByCity(ctx context.Context, city string) ([]models.UserEntity, error)
}

func New(messageStore MessageCreator, userStore UserFinder, log *slog.Logger) *Sender {
	return &Sender{
		messageStore: messageStore,
		userStore:    userStore,
		log:          log,
	}

}

func (s *Sender) Send(message models.MessageConsumer) error {
	usersStore, err := s.userStore.FindByCity(context.Background(), message.City)
	if err != nil {
		// if we don't find any users, we don't return an error
		if errors.Is(err, errorx.ErrNotFound) {
			return nil
		}
		s.log.With(slog.Any("city", message.City)).
			Error("finding users by city", err)
		return err
	}
	users, err := s.transformUsersStoreToUsers(usersStore)
	if err != nil {
		s.log.Error("transforming users store to users", err)
		return err
	}

	usersCh := make(chan *models.User, len(users))
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go s.send(context.Background(), usersCh, message, &wg)
	}

	go writeUsersToChannel(users, usersCh)
	wg.Wait()

	return nil
}

// writeUsersToChannel writes users to the channel
func writeUsersToChannel(users []*models.User, usersCh chan<- *models.User) {
	for _, u := range users {
		usersCh <- u
	}
	close(usersCh)
}

// send sends the message to the users
func (s *Sender) send(ctx context.Context, usersCh <-chan *models.User, message models.MessageConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range usersCh {
		for _, contact := range user.Contacts {
			if !contact.IsActive {
				continue
			}
			newMessage, err := s.transformMessageToStoreModel(message, user.ID, contact)
			newMessage.ID = uuid.New()
			if err != nil {
				s.log.With(slog.Any("message", message), slog.Any("user_id", user.ID)).
					Error("transforming message to store model", err)
				continue
			}

			if err = s.messageStore.Create(ctx, newMessage); err != nil {
				s.log.With(slog.Any("message", newMessage)).
					Error("creating message", err)
				continue
			}
		}
	}
}

func (s *Sender) transformMessageToStoreModel(m models.MessageConsumer, userID uuid.UUID, contact models.Contact) (*models.MessageEntity, error) {
	storeModel := &models.MessageEntity{
		Subject: m.Subject,
		Text:    m.Text,
		Status:  m.Status,
		UserID:  userID,
		Type:    contact.Type,
		Value:   contact.Value,
	}
	return storeModel, nil
}

func (s *Sender) transformUsersStoreToUsers(usersStore []models.UserEntity) ([]*models.User, error) {
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
