package mail_gun

import (
	"os"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"

	"github.com/mailgun/mailgun-go"
)

type ClientMailg struct {
	mg  *mailgun.MailgunImpl
	log logging.Logger
}

func NewEmailMailgClient(log logging.Logger) *ClientMailg {
	apiKey := os.Getenv("EMAIL_API_KEY")
	domain := os.Getenv("EMAIL_DOMAIN")

	mg := mailgun.NewMailgun(domain, apiKey)

	return &ClientMailg{
		mg:  mg,
		log: log,
	}
}

func (c ClientMailg) message(newMessage models.Message, email string) *mailgun.Message {
	return c.mg.NewMessage(
		os.Getenv("MAILGUN_FROM"),
		newMessage.Subject,
		newMessage.Text,
		email,
	)
}

func (c ClientMailg) Send(newMessage models.Message, email string) error {
	message := c.message(newMessage, email)
	_, _, err := c.mg.Send(message)
	if err != nil {
		c.log.Error("cannot send message:", err)
		return err
	}
	return err
}
