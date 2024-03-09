package providers

import (
	"fmt"
	"log/slog"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	"projects/emergency-messages/internal/providers/sms/twil"
)

type Sender interface {
	Send(message models.MessageSend) error
}

type SendManager struct {
	providers map[models.ContactType]Sender
}

func New(l *slog.Logger) *SendManager {
	mailg := mail_gun.NewEmailMailgClient(l)
	twilSMS := twil.NewMobileTwilClient(l)

	suppliers := map[models.ContactType]Sender{
		models.ContactTypeEmail: mailg,
		models.ContactTypeSMS:   twilSMS,
	}

	return &SendManager{
		providers: suppliers,
	}
}

func (sm *SendManager) addProvider(provider Sender, cType models.ContactType) {
	sm.providers[cType] = provider
}

func (sm *SendManager) Send(message models.MessageSend) error {
	provider := sm.providers[message.Type]
	if provider == nil {
		return fmt.Errorf("couldn't find provider by type: %s", message.Type)
	}

	return provider.Send(message)
}
