package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"projects/emergency-messages/internal/models"
	"strconv"
)

const numberOfCSVCells = 7
const semicolon = ';'

type UserService struct {
	userStore User
	log       *slog.Logger
}

type User interface {
	Create(ctx context.Context, user *models.UserEntity) error
	FindByCity(ctx context.Context, city string) ([]models.UserEntity, error)
}

func NewUserService(userStore User, log *slog.Logger) UserService {
	return UserService{
		userStore: userStore,
		log:       log,
	}
}

func (s *UserService) GetByCity(ctx context.Context, city string) ([]models.User, error) {
	if city == "" {
		err := errors.New("city is empty")
		s.log.Error("checking input data", err)
		return nil, err
	}
	usersStore, err := s.userStore.FindByCity(ctx, city)
	if err != nil {
		s.log.Error("finding by city", err)
		return nil, err
	}

	users, err := s.transformStoreModelsByCityToUsers(usersStore)
	if err != nil {
		s.log.Error("transforming store model to user", err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) Upload(csvData io.Reader) ([]*models.UserCreate, error) {
	ctx := context.Background()
	users, err := s.getUsersFromCSV(csvData)
	if err != nil {
		s.log.Error("getting users from csv", err)
		return nil, err
	}

	result := make([]*models.UserCreate, 0, len(users))

	for _, user := range users {
		userStore, err := s.transformUserCreateToStoreModel(user)
		if err != nil {
			s.log.Error("transforming user to store model", err)
			return nil, err
		}
		if err = s.userStore.Create(ctx, userStore); err != nil {
			s.log.Error("creating user", err)
			return nil, err
		}

		userCreated, err := s.transformStoreModelToUser(userStore)
		if err != nil {
			s.log.Error("transforming store model to user", err)
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
		s.log.Error("reading csv", err)
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
		// v[3] is IsMobileActive
		// v[4] is Email
		// v[5] is IsEmailActive
		// v[6] is City
		firstName := v[0]
		lastName := v[1]
		if firstName == "" || lastName == "" {
			continue
		}

		contacts, err := getContacts(v)
		if err != nil {
			continue
		}

		user := &models.UserCreate{
			FirstName: firstName,
			LastName:  lastName,
			Contacts:  contacts,
			City:      v[6],
		}

		users = append(users, user)
	}

	return users, nil
}

func getContacts(v []string) ([]models.Contact, error) {
	// v[2] is MobilePhone
	// v[3] is IsMobileActive
	var contacts []models.Contact
	if mobilePhone := v[2]; mobilePhone != "" {
		isMobileActive, err := strconv.ParseBool(v[3])
		if err != nil {
			return nil, err
		}
		contact := models.Contact{
			Value:    mobilePhone,
			Type:     models.ContactTypeSMS,
			IsActive: isMobileActive,
		}
		contacts = append(contacts, contact)
	}
	// v[4] is Email
	// v[5] is IsEmailActive
	if email := v[4]; email != "" {
		isEmailActive, err := strconv.ParseBool(v[5])
		if err != nil {
			return nil, err
		}
		contact := models.Contact{
			Value:    email,
			Type:     models.ContactTypeEmail,
			IsActive: isEmailActive,
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (s *UserService) transformStoreModelsByCityToUsers(usersStore []models.UserEntity) ([]models.User, error) {
	users := make([]models.User, len(usersStore))
	for _, u := range usersStore {
		user := models.User{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Contacts:  u.Contacts,
			City:      u.City,
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *UserService) transformUserCreateToStoreModel(u *models.UserCreate) (*models.UserEntity, error) {
	return &models.UserEntity{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Contacts:  u.Contacts,
		City:      u.City,
	}, nil
}

func (s *UserService) transformStoreModelToUser(u *models.UserEntity) (*models.UserCreate, error) {
	return &models.UserCreate{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Contacts:  u.Contacts,
		City:      u.City,
	}, nil
}
