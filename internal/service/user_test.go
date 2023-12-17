package service

import (
	"bytes"
	"context"
	"errors"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	"testing"

	mock_store "projects/emergency-messages/internal/store/postgres/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetByCity(t *testing.T) {
	t.Skip()
	t.Run("when have city and service without error then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		city := "Moscow"

		wantReturn := []models.User{
			{
				ID:          "1",
				FirstName:   "Albert",
				LastName:    "Guss",
				MobilePhone: "+8748327432",
				Email:       "al@gmail.com",
				City:        "Vancouver",
			},
		}
		userStore.
			EXPECT().
			FindByCity(ctx, city).
			Return(wantReturn, nil)

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.GetByCity(ctx, city)
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
	t.Run("when city is empty then error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		city := ""

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.GetByCity(ctx, city)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
	t.Run("when have city, but store return error then error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		city := "Sao Paulo"

		userStore.
			EXPECT().
			FindByCity(ctx, city).
			Return(nil, errors.New(""))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.GetByCity(ctx, city)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserService_Upload(t *testing.T) {
	t.Skip()
	t.Run("when all data have then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mock_store.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			ID:          "111",
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(nil)

		data := "FirstName;SecondName;MobilePhone;Email;City\nAlbert;Guss;+8748327432;al@gmail.com;Paris\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
	t.Run("when first name in csv is empty then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mock_store.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			ID:          "0e9876fa-b3f2-4d81-9fc9-182641bfa8b0",
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(nil)

		data := "FirstName;SecondName;MobilePhone;Email;City\n;Smith;+4723746273;ezolda@gmail.com;Berlin\nAlbert;Guss;+8748327432;al@gmail.com;Paris\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
	t.Run("when last name in csv is empty then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mock_store.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			ID:          "1",
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(nil)

		data := "FirstName;SecondName;MobilePhone;Email;City\nEzolda;;+4723746273;ezolda@gmail.com;Berlin\nAlbert;Guss;+8748327432;al@gmail.com;Paris\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
	t.Run("when csv is invalid, there are too few cells then error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mock_store.NewMockUser(ctrl)

		data := "FirstName;SecondName;MobilePhone;Email\nAlbert;Guss;+8748327432;al@gmail.com\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
	t.Run("when csv is valid, but store returns error then error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mock_store.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			ID:          "1",
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(errors.New(""))

		data := "FirstName;SecondName;MobilePhone;Email;City\nAlbert;Guss;+8748327432;al@gmail.com;Paris\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
