package services

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetByCity(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userStore := mock_services.NewMockUser(ctrl)
	userService := NewUserService(userStore, log)

	t.Run("when have city and services then no error", func(t *testing.T) {
		city := "Moscow"
		wantReturn := []models.UserEntity{
			{
				FirstName: "Albert",
				LastName:  "Guss",
				Contacts: []models.Contact{
					{
						Value: "+8748327432",
						Type:  models.ContactTypeSMS,
					},
					{
						Value:    "al@gmail.com",
						Type:     models.ContactTypeEmail,
						IsActive: true,
					},
				},
				City: "Vancouver",
			},
		}
		userStore.
			EXPECT().
			FindByCity(ctx, city).
			Return(wantReturn, nil)

		users, err := userService.GetByCity(ctx, city)
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})

	t.Run("when city is empty then error", func(t *testing.T) {
		users, err := userService.GetByCity(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("when have city, but stores return error then error", func(t *testing.T) {
		city := "Sao Paulo"

		userStore.
			EXPECT().
			FindByCity(ctx, city).
			Return(nil, errors.New(""))

		users, err := userService.GetByCity(ctx, city)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserService_Upload(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userStore := mock_services.NewMockUser(ctrl)
	ctx := context.Background()

	userService := NewUserService(userStore, log)

	t.Run("when all data have then no error", func(t *testing.T) {
		userCreate := &models.UserEntity{
			FirstName: "Robert",
			LastName:  "Smith",
			Contacts: []models.Contact{
				{
					Value: "+48178323",
					Type:  models.ContactTypeSMS,
				},
				{
					Value:    "iaiw3br@gmail.com",
					Type:     models.ContactTypeEmail,
					IsActive: true,
				},
			},
			City: "Saint-Petersburg",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(nil)

		data := "firstName;secondName;MobilePhone;IsMobileActive;Email;IsEmailActive;City\nRobert;Smith;+48178323;false;iaiw3br@gmail.com;true;Saint-Petersburg"
		buf := bytes.NewBuffer([]byte(data))

		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
	t.Run("when first name in csv is empty then no error", func(t *testing.T) {
		userCreate := &models.UserEntity{
			FirstName: "Robert",
			LastName:  "Smith",
			Contacts: []models.Contact{
				{
					Value: "+48178323",
					Type:  models.ContactTypeSMS,
				},
				{
					Value:    "iaiw3br@gmail.com",
					Type:     models.ContactTypeEmail,
					IsActive: true,
				},
			},
			City: "Saint-Petersburg",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(nil)

		data := "firstName;secondName;MobilePhone;IsMobileActive;Email;IsEmailActive;City\nBoris;;+7312787312;false;iaiw3br@gmail.com;true;Saint-Petersburg\nRobert;Smith;+48178323;false;iaiw3br@gmail.com;true;Saint-Petersburg"
		buf := bytes.NewBuffer([]byte(data))

		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
	t.Run("when last name in csv is empty then no error", func(t *testing.T) {

		data := "firstName;secondName;MobilePhone;IsMobileActive;Email;IsEmailActive;City\nBoris;;+7312787312;false;iaiw3br@gmail.com;true;Saint-Petersburg\nRobert;;+48178323;false;iaiw3br@gmail.com;true;Saint-Petersburg"
		buf := bytes.NewBuffer([]byte(data))

		users, err := userService.Upload(buf)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})
	t.Run("when csv is invalid, there are too few cells then error", func(t *testing.T) {
		data := "FirstName;SecondName;MobilePhone;Email\nAlbert;Guss;+8748327432;al@gmail.com\n"
		buf := bytes.NewBuffer([]byte(data))

		users, err := userService.Upload(buf)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
	t.Run("when csv is valid, but stores returns error then error", func(t *testing.T) {
		userCreate := &models.UserEntity{
			FirstName: "Boris",
			LastName:  "Ivanov",
			Contacts: []models.Contact{
				{
					Value: "+7312787312",
					Type:  models.ContactTypeSMS,
				},
				{
					Value:    "iaiw3br@gmail.com",
					Type:     models.ContactTypeEmail,
					IsActive: true,
				},
			},
			City: "Saint-Petersburg",
		}
		userStore.
			EXPECT().
			Create(ctx, userCreate).
			Return(errors.New(""))

		data := "firstName;secondName;MobilePhone;IsMobileActive;Email;IsEmailActive;City\nBoris;Ivanov;+7312787312;false;iaiw3br@gmail.com;true;Saint-Petersburg\nRobert;Smith;+48178323;false;iaiw3br@gmail.com;true;Saint-Petersburg"
		buf := bytes.NewBuffer([]byte(data))

		users, err := userService.Upload(buf)
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
