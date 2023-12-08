package providers

import "github.com/emergency-messages/internal/models"

type Sender interface {
	Send(newMessage models.Message, email string) error
}
