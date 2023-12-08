package models

import "errors"

type Template struct {
	ID      uint64 `json:"id"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type TemplateUpdate struct {
	ID      uint64 `json:"id"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func (t TemplateUpdate) Validate() error {
	if t.ID == uint64(0) {
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
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func (t TemplateCreate) Validate() error {
	if t.Subject == "" {
		return errors.New("subject is empty")
	}
	if t.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}
