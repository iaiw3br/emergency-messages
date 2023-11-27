package store

import (
	"context"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *pgx.Conn
}

// type User interface {
// 	CreateMany(ctx context.Context, users []models.User) error
// }

func NewUserStore(db *pgx.Conn) UserStore {
	return UserStore{
		db: db,
	}
}

func (u UserStore) CreateMany(ctx context.Context, users []models.User) error {
	sql := `
		INSERT INTO users (first_name, last_name, mobile_phone, email) 
		VALUES ($1, $2, $3, $4)
	`
	for _, v := range users {
		_, err := u.db.Exec(ctx, sql, v.FirstName, v.LastName, v.MobilePhone, v.Email)
		if err != nil {
			return err
		}
	}
	return nil
}
