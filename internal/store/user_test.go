package store

import (
	"context"
	"github.com/emergency-messages/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStore_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDatabase(t)

	user := &models.UserCreate{
		FirstName:   "David",
		LastName:    "Smith",
		MobilePhone: "+87467328423",
		Email:       "david@gmail.com",
		City:        "Moscow",
	}

	userStore := NewUserStore(db)
	id, err := userStore.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestUserStore_FindByCity(t *testing.T) {
	ctx := context.Background()
	db := setupTestDatabase(t)

	user := &models.UserCreate{
		FirstName:   "David",
		LastName:    "Smith",
		MobilePhone: "+87467328423",
		Email:       "david@gmail.com",
		City:        "Moscow",
	}
	userStore := NewUserStore(db)
	id, err := userStore.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	users, err := userStore.FindByCity(ctx, user.City)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
}
