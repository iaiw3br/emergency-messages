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

type ReceiverService struct {
	receiverStore ReceiverStore
	log           *slog.Logger
}

type ReceiverStore interface {
	Create(ctx context.Context, receiver *models.ReceiverEntity) error
	FindByCity(ctx context.Context, city string) ([]models.ReceiverEntity, error)
}

func NewReceiverService(receiverStore ReceiverStore, log *slog.Logger) *ReceiverService {
	return &ReceiverService{
		receiverStore: receiverStore,
		log:           log,
	}
}

func (s *ReceiverService) FindByCity(ctx context.Context, city string) ([]models.Receiver, error) {
	if city == "" {
		err := errors.New("city is empty")
		s.log.Error("checking input queue", err)
		return nil, err
	}
	receiverStore, err := s.receiverStore.FindByCity(ctx, city)
	if err != nil {
		s.log.Error("finding by city", err)
		return nil, err
	}

	receivers, err := s.transformStoreModelsByCityToReceivers(receiverStore)
	if err != nil {
		s.log.Error("transforming store model to ReceiverService", err)
		return nil, err
	}

	return receivers, nil
}

func (s *ReceiverService) Upload(csvData io.Reader) ([]*models.ReceiverCreate, error) {
	ctx := context.Background()
	receivers, err := s.getReceiversFromCSV(csvData)
	if err != nil {
		s.log.Error("getting receivers from csv", err)
		return nil, err
	}

	result := make([]*models.ReceiverCreate, 0, len(receivers))

	for _, receiver := range receivers {
		receiverStoreModel, err := s.transformReceiverCreateToStoreModel(receiver)
		if err != nil {
			s.log.Error("transforming ReceiverService to store model", err)
			return nil, err
		}
		if err = s.receiverStore.Create(ctx, receiverStoreModel); err != nil {
			s.log.Error("creating ReceiverService", err)
			return nil, err
		}

		receiverCreated, err := s.transformStoreModelToReceiver(receiverStoreModel)
		if err != nil {
			s.log.Error("transforming store model to ReceiverService", err)
			return nil, err
		}
		result = append(result, receiverCreated)
	}
	return result, nil
}

func (s *ReceiverService) getReceiversFromCSV(csvData io.Reader) ([]*models.ReceiverCreate, error) {
	csvReader := csv.NewReader(csvData)
	csvReader.Comma = semicolon

	records, err := csvReader.ReadAll()
	if err != nil {
		s.log.Error("reading csv", err)
		return nil, err
	}

	receivers := make([]*models.ReceiverCreate, 0, len(records))
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

		receiver := &models.ReceiverCreate{
			FirstName: firstName,
			LastName:  lastName,
			Contacts:  contacts,
			City:      v[6],
		}

		receivers = append(receivers, receiver)
	}

	return receivers, nil
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

func (s *ReceiverService) transformStoreModelsByCityToReceivers(receiverStore []models.ReceiverEntity) ([]models.Receiver, error) {
	receivers := make([]models.Receiver, len(receiverStore))
	for _, u := range receiverStore {
		receiver := models.Receiver{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Contacts:  u.Contacts,
			City:      u.City,
		}
		receivers = append(receivers, receiver)
	}
	return receivers, nil
}

func (s *ReceiverService) transformReceiverCreateToStoreModel(u *models.ReceiverCreate) (*models.ReceiverEntity, error) {
	return &models.ReceiverEntity{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Contacts:  u.Contacts,
		City:      u.City,
	}, nil
}

func (s *ReceiverService) transformStoreModelToReceiver(u *models.ReceiverEntity) (*models.ReceiverCreate, error) {
	return &models.ReceiverCreate{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Contacts:  u.Contacts,
		City:      u.City,
	}, nil
}
