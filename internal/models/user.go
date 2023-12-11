package models

import "time"

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MobilePhone string    `json:"mobile_phone"`
	Email       string    `json:"email"`
	City        string    `json:"city"`
	Created     time.Time `json:"created"`
}

type UserCreate struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MobilePhone string    `json:"mobile_phone"`
	Email       string    `json:"email"`
	City        string    `json:"city"`
	Created     time.Time `json:"created"`
}
