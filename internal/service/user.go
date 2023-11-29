package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/store"
	"io"
)

const numberOfCSVCells = 4
const semicolon = ';'

type UserService struct {
	userStore store.User
	log       logging.Logger
}

func NewUserService(userStore store.User, log logging.Logger) UserService {
	return UserService{
		userStore: userStore,
		log:       log,
	}
}

func (u UserService) Upload(csvData io.Reader) error {
	ctx := context.Background()
	users, err := u.getUsersFromCSV(csvData)
	if err != nil {
		u.log.Error(err)
		return err
	}
	return u.userStore.CreateMany(ctx, users)
}

func (u UserService) getUsersFromCSV(csvData io.Reader) ([]models.User, error) {
	csvReader := csv.NewReader(csvData)
	csvReader.Comma = semicolon

	records, err := csvReader.ReadAll()
	if err != nil {
		u.log.Error(err)
		return nil, err
	}

	users := make([]models.User, 0, len(records))
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

		firstName := v[0]
		lastName := v[1]
		if firstName == "" || lastName == "" {
			continue
		}

		user := models.User{
			FirstName:   firstName,
			LastName:    lastName,
			MobilePhone: v[2],
			Email:       v[3],
		}

		users = append(users, user)
	}

	return users, nil
}
