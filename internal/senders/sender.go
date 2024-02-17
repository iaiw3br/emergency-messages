// Package senders
// Implements logic to send messages to clients
package senders

import (
	"context"
	"projects/emergency-messages/internal/models"
	"sync"
)

type Sender struct{}

func New() *Sender {
	return &Sender{}
}

func (s *Sender) Send() {
	// do something
}

func (s *Sender) send(ctx context.Context, usersCh <-chan *models.User, newMessage models.Message, wg *sync.WaitGroup) {
	// defer wg.Done()
	// for user := range usersCh {
	// 	for _, contact := range user.Contacts {
	// 		if !contact.IsActive {
	// 			continue
	// 		}
	// 		newMessage.UserID = user.ID
	//
	// 		storeModel, err := s.transformMessageToStoreModel(newMessage)
	// 		if err != nil {
	// 			s.log.Error("transforming message to store model", err)
	// 			continue
	// 		}
	//
	// 		if err = s.messageStore.Create(ctx, storeModel); err != nil {
	// 			s.log.Error("creating message", err)
	// 			continue
	// 		}
	//
	// 		if err = s.sender.Send(newMessage, contact); err != nil {
	// 			s.log.Error("sending message", err)
	// 			continue
	// 		}
	//
	// 		if err = s.messageStore.UpdateStatus(ctx, storeModel.ID, models.Delivered); err != nil {
	// 			s.log.Error("updating message", err)
	// 			continue
	// 		}
	// 	}
	// }
}

//
// func (s *MessageService) transformMessageToStoreModel(m models.Message) (*models.MessageEntity, error) {
// 	storeModel := &models.MessageEntity{
// 		Subject: m.Subject,
// 		Text:    m.Text,
// 		Status:  m.Status,
// 		UserID:  m.UserID,
// 	}
// 	return storeModel, nil
// }
//
// func (s *MessageService) transformUsersStoreToUsers(usersStore []models.UserEntity) ([]*models.User, error) {
// 	users := make([]*models.User, 0, len(usersStore))
// 	for _, u := range usersStore {
// 		user := &models.User{
// 			ID:        u.ID,
// 			FirstName: u.FirstName,
// 			LastName:  u.LastName,
// 			Contacts:  u.Contacts,
// 			City:      u.City,
// 		}
// 		users = append(users, user)
// 	}
// 	return users, nil
// }
//
// func sendUsersToUsersChannel(users []*models.User, usersCh chan<- *models.User) {
// 	for _, u := range users {
// 		usersCh <- u
// 	}
// 	close(usersCh)
// }

// func readFromQueue() {
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:53517",
// 		"group.id":          "myGroup",
// 		"auto.offset.reset": "earliest",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	topic := "send"
// 	if err = c.SubscribeTopics([]string{topic}, nil); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
//
// 	// A signal handler or similar could be used to set this to false to break the loop.
// 	run := true
// 	for run {
// 		select {
// 		case sig := <-sigchan:
// 			fmt.Printf("Caught signal %v: terminating\n", sig)
// 			run = false
// 		default:
// 			msg, err := c.ReadMessage(time.Second)
// 			if err == nil {
// 				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
// 			} else if !err.(kafka.Error).IsTimeout() {
// 				// The client will automatically try to recover from all errors.
// 				// Timeout is not considered an error because it is raised by
// 				// ReadMessage in absence of messages.
// 				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 			}
// 		}
// 	}
// 	c.Close()
// }
