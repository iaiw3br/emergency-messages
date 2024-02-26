package postgres

import (
	"context"
	"fmt"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services"

	"github.com/uptrace/bun"
)

type receiverStore struct {
	db *bun.DB
}

func NewReceiverStore(db *bun.DB) services.ReceiverStore {
	return &receiverStore{
		db: db,
	}
}

// Create creates the struct of a receiver in the database.
// It takes in a context, the new struct of the receiver.
// It returns an error if the create operation fails.
func (s *receiverStore) Create(ctx context.Context, u *models.ReceiverEntity) error {
	_, err := s.db.
		NewInsert().
		Model(u).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("creating receiver: couldn't create: %v. Error: %w", u, err)
	}
	return nil
}

// FindByCity retrieves receivers from the database by city.
// It takes in a context and the city of the receiver.
// It returns receivers and an error if the retrieval operation fails.
func (s *receiverStore) FindByCity(ctx context.Context, city string) ([]models.ReceiverEntity, error) {
	entities := make([]models.ReceiverEntity, 0)

	err := s.db.
		NewSelect().
		Model(&entities).
		Where("city = ?", city).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding receivers by city: couldn't find receivers by city: %s. Error: %w", city, err)
	}

	if len(entities) == 0 {
		return nil, errorx.ErrNotFound
	}

	return entities, nil
}
