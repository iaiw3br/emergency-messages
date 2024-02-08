package twil

import (
	"os"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ClientTwilSMS struct {
	log logging.Logger
}

func (c *ClientTwilSMS) Send(newMessage models.Message, phone string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("MOBILE_TWIL_ACCOUNT_SID"),
		Password: os.Getenv("MOBILE_TWIL_AUTH_TOKEN"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(os.Getenv("MOBILE_PHONE_EMERGENCY_SERVICE"))
	params.SetBody(newMessage.Text)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		c.log.Errorf("mobile twil couldn't send message: %v, phone: %s. Error: %v", newMessage, phone, err)
		return err
	}
	return nil
}

func NewMobileTwilClient(log logging.Logger) *ClientTwilSMS {
	return &ClientTwilSMS{
		log: log,
	}
}
