package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/google/uuid"
	"io"
	"time"
)

const numberOfCSVCells = 5
const semicolon = ';'

type UserService struct {
	userStore User
	log       logging.Logger
}

type User interface {
	Create(ctx context.Context, user *models.UserCreate) error
	FindByCity(ctx context.Context, city string) ([]models.User, error)
}

func NewUserService(userStore User, log logging.Logger) UserService {
	return UserService{
		userStore: userStore,
		log:       log,
	}
}

func (u UserService) GetByCity(ctx context.Context, city string) ([]models.User, error) {
	if city == "" {
		err := errors.New("city is empty")
		u.log.Error(err)
		return nil, err
	}
	users, err := u.userStore.FindByCity(ctx, city)
	if err != nil {
		u.log.Error(err)
		return nil, err
	}

	return users, nil
}

func (u UserService) Upload(csvData io.Reader) ([]*models.UserCreate, error) {
	ctx := context.Background()
	users, err := u.getUsersFromCSV(csvData)
	if err != nil {
		u.log.Error(err)
		return nil, err
	}

	result := make([]*models.UserCreate, 0, len(users))
	now := time.Now()

	for _, user := range users {
		user.ID = uuid.New().String()
		user.Created = now
		if err = u.userStore.Create(ctx, user); err != nil {
			u.log.Error(err)
			return nil, err
		}

		userCreated := &models.UserCreate{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			MobilePhone: user.MobilePhone,
			Email:       user.Email,
			City:        user.City,
		}
		result = append(result, userCreated)
	}
	return result, nil
}

func (u UserService) getUsersFromCSV(csvData io.Reader) ([]*models.UserCreate, error) {
	csvReader := csv.NewReader(csvData)
	csvReader.Comma = semicolon

	records, err := csvReader.ReadAll()
	if err != nil {
		u.log.Error(err)
		return nil, err
	}

	users := make([]*models.UserCreate, 0, len(records))
	for i, v := range records {
		if len(v) != numberOfCSVCells {
			return nil, errors.New(fmt.Sprintf("invalid csv file, expect %d cells, have %d cells", numberOfCSVCells, len(v)))
		}

		// the head of the file
		if i == 0 {
			continue
		}
		// v[0] is firstName
		// v[1] is secondName
		// v[2] is MobilePhone
		// v[3] is Email
		// v[4] is City

		firstName := v[0]
		lastName := v[1]
		if firstName == "" || lastName == "" {
			continue
		}

		user := &models.UserCreate{
			FirstName:   firstName,
			LastName:    lastName,
			MobilePhone: v[2],
			Email:       v[3],
			City:        v[4],
		}

		users = append(users, user)
	}

	return users, nil
}
