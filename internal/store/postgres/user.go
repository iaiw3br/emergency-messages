package postgres

import (
	"context"
	"fmt"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/service"

	"github.com/uptrace/bun"
)

type userStore struct {
	db *bun.DB
}

func NewUserStore(db *bun.DB) service.User {
	return &userStore{
		db: db,
	}
}

// Create creates the struct of a user in the database.
// It takes in a context, the new struct of the user.
// It returns an error if the create operation fails.
func (s *userStore) Create(ctx context.Context, u *models.UserEntity) error {
	_, err := s.db.
		NewInsert().
		Model(u).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("creating user: couldn't create: %v. Error: %w", u, err)
	}
	return nil
}

// FindByCity retrieves users from the database by city.
// It takes in a context and the city of the user.
// It returns users and an error if the retrieval operation fails.
func (s *userStore) FindByCity(ctx context.Context, city string) ([]models.UserEntity, error) {
	entities := make([]models.UserEntity, 0)

	err := s.db.
		NewSelect().
		Model(&entities).
		Where("city = ?", city).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding users by city: couldn't find users by city: %s. Error: %w", city, err)
	}

	if len(entities) == 0 {
		return nil, fmt.Errorf("finding users by city: couldn't find users by city: %s", city)
	}

	return entities, nil
}
