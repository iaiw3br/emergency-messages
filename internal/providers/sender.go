package providers

import (
	"fmt"
	"projects/emergency-messages/internal/models"
)

type Sender interface {
	Send(message models.Message, where string) error
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

func (sm *SendManager) Send(message models.Message, contact models.Contact) error {
	provider := sm.providers[contact.Type]
	if provider == nil {
		return fmt.Errorf("couldn't find provider by type: %s", contact.Type)
	}

	return provider.Send(message, contact.Value)
}
