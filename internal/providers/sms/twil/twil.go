package twil

import (
	"log/slog"
	"os"
	"projects/emergency-messages/internal/models"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ClientTwilSMS struct {
	log *slog.Logger
}

func (c *ClientTwilSMS) Send(newMessage models.MessageSend) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("MOBILE_TWIL_ACCOUNT_SID"),
		Password: os.Getenv("MOBILE_TWIL_AUTH_TOKEN"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(newMessage.Value)
	params.SetFrom(os.Getenv("MOBILE_PHONE_EMERGENCY_SERVICE"))
	params.SetBody(newMessage.Text)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		c.log.With(
			slog.Any("message", newMessage),
			slog.String("phone", newMessage.Value)).
			Error("sending twil message", err)
		return err
	}
	return nil
}

func NewMobileTwilClient(log *slog.Logger) *ClientTwilSMS {
	return &ClientTwilSMS{
		log: log,
	}
}
