package models

type User struct {
	ID          uint64 `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MobilePhone string `json:"mobile_phone"`
	Email       string `json:"email"`
	City        string `json:"city"`
}
