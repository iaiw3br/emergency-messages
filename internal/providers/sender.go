package providers

import (
	"fmt"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	"projects/emergency-messages/internal/providers/sms/twil"
)

type Client struct {
	mailGun *mail_gun.ClientMailg
	smsTwil *twil.ClientTwilSMS
}

func NewClient(log logging.Logger) Client {
	mailg := mail_gun.NewEmailMailgClient(log)
	smsTwil := twil.NewMobileTwilClient(log)

	return Client{
		mailGun: mailg,
		smsTwil: smsTwil,
	}
}

func (c *Client) Send(message models.Message, contact models.Contact) error {
	switch contact.Type {
	case models.ContactTypeEmail:
		return c.mailGun.Send(message, contact.Value)
	case models.ContactTypeSMS:
		return c.smsTwil.Send(message, contact.Value)
	default:
		return fmt.Errorf("bad contact type: %s", contact.Type)
	}
}
