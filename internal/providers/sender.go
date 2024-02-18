package providers

import (
	"fmt"
	"projects/emergency-messages/internal/models"
)

type Sender interface {
	Send(message models.MessageSend) error
}

type SendManager struct {
	providers map[models.ContactType]Sender
}

func New() *SendManager {
	return &SendManager{
		providers: make(map[models.ContactType]Sender),
	}
}

func (sm *SendManager) AddProvider(provider Sender, cType models.ContactType) {
	sm.providers[cType] = provider
}

func (sm *SendManager) Send(message models.MessageSend) error {
	provider := sm.providers[message.Type]
	if provider == nil {
		return fmt.Errorf("couldn't find provider by type: %s", message.Type)
	}

	return provider.Send(message)
}
