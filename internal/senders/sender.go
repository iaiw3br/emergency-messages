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
	messageStore  MessageCreator
	receiverStore ReceiverFinder
	log           *slog.Logger
}

type MessageCreator interface {
	Create(ctx context.Context, m *models.MessageEntity) error
}

type ReceiverFinder interface {
	FindByCity(ctx context.Context, city string) ([]models.ReceiverEntity, error)
}

func New(messageStore MessageCreator, receiverStore ReceiverFinder, log *slog.Logger) *Sender {
	return &Sender{
		messageStore:  messageStore,
		receiverStore: receiverStore,
		log:           log,
	}

}

func (s *Sender) Send(message models.MessageConsumer) error {
	receiverStore, err := s.receiverStore.FindByCity(context.Background(), message.City)
	if err != nil {
		// if we don't find any receivers, we don't return an error
		if errors.Is(err, errorx.ErrNotFound) {
			return nil
		}
		s.log.With(slog.Any("city", message.City)).
			Error("finding receivers by city", err)
		return err
	}
	receivers, err := s.transformReceiversStoreToReceivers(receiverStore)
	if err != nil {
		s.log.Error("transforming receivers store to receivers", err)
		return err
	}

	receiversCh := make(chan *models.Receiver, len(receivers))
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go s.send(context.Background(), receiversCh, message, &wg)
	}

	go writeReceiversToChannel(receivers, receiversCh)
	wg.Wait()

	return nil
}

// writeReceiversToChannel writes receivers to the channel
func writeReceiversToChannel(receivers []*models.Receiver, receiversCh chan<- *models.Receiver) {
	for _, u := range receivers {
		receiversCh <- u
	}
	close(receiversCh)
}

// send sends the message to the receivers
func (s *Sender) send(ctx context.Context, receiversCh <-chan *models.Receiver, message models.MessageConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	for receiver := range receiversCh {
		for _, contact := range receiver.Contacts {
			if !contact.IsActive {
				continue
			}
			newMessage, err := s.transformMessageToStoreModel(message, receiver.ID, contact)
			newMessage.ID = uuid.New()
			if err != nil {
				s.log.With(slog.Any("message", message), slog.Any("receiver_id", receiver.ID)).
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

func (s *Sender) transformMessageToStoreModel(m models.MessageConsumer, receiverID uuid.UUID, contact models.Contact) (*models.MessageEntity, error) {
	storeModel := &models.MessageEntity{
		Subject:    m.Subject,
		Text:       m.Text,
		Status:     m.Status,
		ReceiverID: receiverID,
		Type:       contact.Type,
		Value:      contact.Value,
	}
	return storeModel, nil
}

func (s *Sender) transformReceiversStoreToReceivers(receiversStore []models.ReceiverEntity) ([]*models.Receiver, error) {
	results := make([]*models.Receiver, 0, len(receiversStore))
	for _, u := range receiversStore {
		receiver := &models.Receiver{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Contacts:  u.Contacts,
			City:      u.City,
		}
		results = append(results, receiver)
	}
	return results, nil
}
