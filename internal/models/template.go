package models

type Template struct {
	ID      uint64 `json:"id"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}
