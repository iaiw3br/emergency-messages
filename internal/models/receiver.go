package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Receiver struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	Contacts  []Contact `json:"contacts"`
}

type Contact struct {
	Value    string      `json:"value"`
	Type     ContactType `json:"type"`
	IsActive bool        `json:"is_active"`
}

func (c *Contact) IsActiveMobilePhone() bool {
	return c.IsActive && c.Type == ContactTypeSMS
}

func (c *Contact) IsActiveEmail() bool {
	return c.IsActive && c.Type == ContactTypeEmail
}

type ContactType string

const (
	ContactTypeEmail ContactType = "email"
	ContactTypeSMS   ContactType = "sms"
)

type ReceiverCreate struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Contacts  []Contact `json:"contacts"`
	City      string    `json:"city"`
}

type ReceiverEntity struct {
	bun.BaseModel `bun:"table:receivers,alias:u"`
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	FirstName     string    `bun:"first_name,notnull"`
	LastName      string    `bun:"last_name,notnull"`
	Contacts      []Contact `bun:"contacts,notnull"`
	City          string    `bun:"city,notnull"`
}

type ReceiverSend struct {
	ID uuid.UUID `json:"id"`
}
