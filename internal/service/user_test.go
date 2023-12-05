package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	mockstore "github.com/emergency-messages/internal/store/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserService_GetByCity(t *testing.T) {
	t.Run("when have city and service without error then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		city := "Moscow"

		wantReturn := []models.User{
			{
				ID:          1,
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
	t.Run("when all data have then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(uint64(0), nil)

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

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(uint64(0), nil)

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

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(uint64(0), nil)

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

		userStore := mockstore.NewMockUser(ctrl)

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

		userStore := mockstore.NewMockUser(ctrl)
		ctx := context.Background()
		userCreate := &models.UserCreate{
			FirstName:   "Albert",
			LastName:    "Guss",
			MobilePhone: "+8748327432",
			Email:       "al@gmail.com",
			City:        "Paris",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(uint64(0), errors.New(""))

		data := "FirstName;SecondName;MobilePhone;Email;City\nAlbert;Guss;+8748327432;al@gmail.com;Paris\n"
		buf := bytes.NewBuffer([]byte(data))

		log := logging.New()
		userService := NewUserService(userStore, log)
		users, err := userService.Upload(buf)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
