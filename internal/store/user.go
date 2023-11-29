package store

import (
	"context"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type User struct {
	db *pgx.Conn
}

// type User interface {
// 	CreateMany(ctx context.Context, users []models.User) error
// }

func NewUserStore(db *pgx.Conn) User {
	return User{
		db: db,
	}
}

func (u User) CreateMany(ctx context.Context, users []models.User) error {
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

// FindByCity find all users by city
func (u User) FindByCity(ctx context.Context, city string) ([]models.User, error) {
	sql := `
		SELECT id, first_name, last_name, email, mobile_phone
		FROM users
		WHERE city = $1;
	`

	rows, err := u.db.Query(ctx, sql, city)
	// FIXME: add not found
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.MobilePhone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
