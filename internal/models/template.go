package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Template struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type TemplateUpdate struct {
	ID      string
	Subject string
	Text    string
}

func (t TemplateUpdate) Validate() error {
	if t.ID == "" {
		return errors.New("id is empty")
	}
	if t.Subject == "" {
		return errors.New("subject is empty")
	}
	if t.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}

type TemplateCreate struct {
	ID      string
	Subject string
	Text    string
}

func (t *TemplateCreate) Validate() error {
	if t.Subject == "" {
		return errors.New("subject is empty")
	}
	if t.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}

type TemplateEntity struct {
	bun.BaseModel `bun:"table:templates,alias:t"`
	ID            uuid.UUID  `bun:"type:uuid,default:uuid_generate_v4()"`
	Subject       string     `bun:"subject,notnull"`
	Text          string     `bun:"text,notnull"`
	DeletedAt     *time.Time `bun:",nullzero"`
}
