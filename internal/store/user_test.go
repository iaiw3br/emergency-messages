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

	user := models.User{
		FirstName:   "David",
		LastName:    "Smith",
		MobilePhone: "+87467328423",
		Email:       "david@gmail.com",
		City:        "Moscow",
	}

	userStore := NewUserStore(db)
	userCreated, err := userStore.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, userCreated.FirstName)
	assert.Equal(t, user.LastName, userCreated.LastName)
	assert.Equal(t, user.MobilePhone, userCreated.MobilePhone)
	assert.Equal(t, user.Email, userCreated.Email)
	assert.Equal(t, user.City, userCreated.City)
	assert.NotNil(t, userCreated.ID)
}

func TestUserStore_FindByCity(t *testing.T) {
	ctx := context.Background()
	db := setupTestDatabase(t)

	user := models.User{
		FirstName:   "David",
		LastName:    "Smith",
		MobilePhone: "+87467328423",
		Email:       "david@gmail.com",
		City:        "Moscow",
	}
	userStore := NewUserStore(db)
	userCreated, err := userStore.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, userCreated.FirstName)
	assert.Equal(t, user.LastName, userCreated.LastName)
	assert.Equal(t, user.MobilePhone, userCreated.MobilePhone)
	assert.Equal(t, user.Email, userCreated.Email)
	assert.Equal(t, user.City, userCreated.City)
	assert.NotNil(t, userCreated.ID)

	users, err := userStore.FindByCity(ctx, userCreated.City)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
}
