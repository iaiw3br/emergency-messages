package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/service"
	"github.com/uptrace/bun"
	"time"
)

type userStore struct {
	db *bun.DB
}

type entityUser struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string    `bun:"id,type:uuid"`
	FirstName     string    `bun:"first_name,notnull"`
	LastName      string    `bun:"last_name,notnull"`
	MobilePhone   string    `bun:"mobile_phone"`
	Email         string    `bun:"email"`
	City          string    `bun:"city,notnull"`
	Created       time.Time `json:"created"`
}

func NewUserStore(db *bun.DB) service.User {
	return &userStore{
		db: db,
	}
}

func (s *userStore) Create(ctx context.Context, user *models.UserCreate) error {
	entity := entityUser{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MobilePhone: user.MobilePhone,
		Email:       user.Email,
		City:        user.City,
	}

	_, err := s.db.NewInsert().
		Model(&entity).
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

// FindByCity find all users by city
func (s *userStore) FindByCity(ctx context.Context, city string) ([]models.User, error) {
	entities := make([]entityUser, 0)

	_, err := s.db.NewSelect().
		Model(&entities).
		Where("city = ?", city).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, errors.New(fmt.Sprintf("users were not found by city: %s", city))
	}

	result := make([]models.User, 0, len(entities))
	for _, u := range entities {
		user := models.User{
			ID:          u.ID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			MobilePhone: u.MobilePhone,
			Email:       u.Email,
			City:        u.City,
			Created:     u.Created,
		}
		result = append(result, user)
	}

	return result, nil
}
