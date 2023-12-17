package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"io"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
)

const numberOfCSVCells = 5
const semicolon = ';'

type UserService struct {
	userStore User
	log       logging.Logger
}

type UserEntity struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            uuid.UUID `bun:"type:uuid,primarykey"`
	FirstName     string    `bun:"first_name,notnull"`
	LastName      string    `bun:"last_name,notnull"`
	MobilePhone   string    `bun:"mobile_phone"`
	Email         string    `bun:"email"`
	City          string    `bun:"city,notnull"`
}

type User interface {
	Create(ctx context.Context, user *UserEntity) error
	FindByCity(ctx context.Context, city string) ([]UserEntity, error)
}

func NewUserService(userStore User, log logging.Logger) UserService {
	return UserService{
		userStore: userStore,
		log:       log,
	}
}

func (s *UserService) GetByCity(ctx context.Context, city string) ([]models.User, error) {
	if city == "" {
		err := errors.New("city is empty")
		s.log.Error(err)
		return nil, err
	}
	usersStore, err := s.userStore.FindByCity(ctx, city)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	users, err := s.transformStoreModelsByCityToUsers(usersStore)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) Upload(csvData io.Reader) ([]*models.UserCreate, error) {
	ctx := context.Background()
	users, err := s.getUsersFromCSV(csvData)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	result := make([]*models.UserCreate, 0, len(users))

	for _, user := range users {
		user.ID = uuid.New()
		userStore, err := s.transformUserCreateToStoreModel(user)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		if err = s.userStore.Create(ctx, userStore); err != nil {
			s.log.Error(err)
			return nil, err
		}

		userCreated, err := s.transformStoreModelToUser(userStore)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		result = append(result, userCreated)
	}
	return result, nil
}

func (s *UserService) getUsersFromCSV(csvData io.Reader) ([]*models.UserCreate, error) {
	csvReader := csv.NewReader(csvData)
	csvReader.Comma = semicolon

	records, err := csvReader.ReadAll()
	if err != nil {
		s.log.Error(err)
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

func (s *UserService) transformStoreModelsByCityToUsers(usersStore []UserEntity) ([]models.User, error) {
	users := make([]models.User, len(usersStore))
	for _, u := range usersStore {
		user := models.User{
			ID:          u.ID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			MobilePhone: u.MobilePhone,
			Email:       u.Email,
			City:        u.City,
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *UserService) transformUserCreateToStoreModel(u *models.UserCreate) (*UserEntity, error) {
	return &UserEntity{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		MobilePhone: u.MobilePhone,
		Email:       u.Email,
		City:        u.City,
	}, nil
}

func (s *UserService) transformStoreModelToUser(u *UserEntity) (*models.UserCreate, error) {
	return &models.UserCreate{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		MobilePhone: u.MobilePhone,
		Email:       u.Email,
		City:        u.City,
	}, nil
}
