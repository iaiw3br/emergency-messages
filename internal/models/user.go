package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

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

type UserEntity struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	FirstName     string    `bun:"first_name,notnull"`
	LastName      string    `bun:"last_name,notnull"`
	MobilePhone   string    `bun:"mobile_phone"`
	Email         string    `bun:"email"`
	City          string    `bun:"city,notnull"`
}
