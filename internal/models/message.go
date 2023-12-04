package models

type Message struct {
	ID      uint64        `json:"id"`
	Subject string        `json:"subject"`
	Text    string        `json:"text"`
	Status  MessageStatus `json:"status"`
	UserID  uint64        `json:"user_id"`
}

type MessageStatus string

const (
	Created   MessageStatus = "created"
	Delivered MessageStatus = "delivered"
)

// Create update status to Created
func (m *Message) Create() {
	m.Status = Created
}

// Deliver update status to Delivered
func (m *Message) Deliver() {
	m.Status = Delivered
}

type CreateMessage struct {
	TemplateID int    `json:"template_id"`
	City       string `json:"city"`
	Strength   string `json:"strength"`
}
