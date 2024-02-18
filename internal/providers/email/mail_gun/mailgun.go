package mail_gun

import (
	"log/slog"
	"os"
	"projects/emergency-messages/internal/models"

	"github.com/mailgun/mailgun-go"
)

type ClientMailg struct {
	mg  *mailgun.MailgunImpl
	log *slog.Logger
}

func NewEmailMailgClient(log *slog.Logger) *ClientMailg {
	apiKey := os.Getenv("EMAIL_API_KEY")
	domain := os.Getenv("EMAIL_DOMAIN")

	mg := mailgun.NewMailgun(domain, apiKey)

	return &ClientMailg{
		mg:  mg,
		log: log,
	}
}

func (c ClientMailg) message(newMessage models.MessageSend) *mailgun.Message {
	return c.mg.NewMessage(
		os.Getenv("MAILGUN_FROM"),
		newMessage.Subject,
		newMessage.Text,
		newMessage.Value,
	)
}

func (c ClientMailg) Send(newMessage models.MessageSend) error {
	message := c.message(newMessage)
	_, _, err := c.mg.Send(message)
	if err != nil {
		c.log.Error("cannot sending message:", err)
		return err
	}
	return err
}
