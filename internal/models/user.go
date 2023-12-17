package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MobilePhone string    `json:"mobile_phone"`
	Email       string    `json:"email"`
	City        string    `json:"city"`
}

type UserCreate struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MobilePhone string    `json:"mobile_phone"`
	Email       string    `json:"email"`
	City        string    `json:"city"`
}
