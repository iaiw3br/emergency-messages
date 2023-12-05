package store

import (
	"context"
	"github.com/emergency-messages/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *pgx.Conn
}

type User interface {
	Create(ctx context.Context, user *models.UserCreate) (uint64, error)
	FindByCity(ctx context.Context, city string) ([]models.User, error)
}

func NewUserStore(db *pgx.Conn) User {
	return UserStore{
		db: db,
	}
}

func (u UserStore) Create(ctx context.Context, user *models.UserCreate) (uint64, error) {
	sql := `
		INSERT INTO users (first_name, last_name, mobile_phone, email, city) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	row := u.db.QueryRow(ctx, sql, user.FirstName, user.LastName, user.MobilePhone, user.Email, user.City)
	var result uint64

	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// FindByCity find all users by city
func (u UserStore) FindByCity(ctx context.Context, city string) ([]models.User, error) {
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
