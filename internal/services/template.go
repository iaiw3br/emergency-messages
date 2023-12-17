package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
)

type TemplateService struct {
	templateStore TemplateStore
	log           logging.Logger
}

type TemplateStore interface {
	Create(ctx context.Context, t *models.TemplateEntity) error
	Update(ctx context.Context, t *models.TemplateEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateEntity, error)
}

func NewTemplate(templateStore TemplateStore, log logging.Logger) TemplateService {
	return TemplateService{
		templateStore: templateStore,
		log:           log,
	}
}

func (s *TemplateService) Create(ctx context.Context, template *models.TemplateCreate) error {
	if err := template.Validate(); err != nil {
		s.log.Error(err)
		return err
	}

	storeModel, err := s.transformTemplateCreateToStoreModel(template)
	if err != nil {
		s.log.Error(err)
		return err
	}

	if err = s.templateStore.Create(ctx, storeModel); err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *TemplateService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("deleting template: id is empty")
	}
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("deleting template: couldn't parse id: %s to UUID. Error: %w", id, err)
	}
	return s.templateStore.Delete(ctx, uuidValue)
}

func (s *TemplateService) Update(ctx context.Context, template *models.TemplateUpdate) error {
	if err := template.Validate(); err != nil {
		s.log.Error(err)
		return err
	}
	storeModel, err := s.transformTemplateUpdateToStoreModel(template)
	if err != nil {
		s.log.Error(err)
		return err
	}
	return s.templateStore.Update(ctx, storeModel)
}

func (s *TemplateService) transformTemplateCreateToStoreModel(t *models.TemplateCreate) (*models.TemplateEntity, error) {
	storeModel := &models.TemplateEntity{
		Subject: t.Subject,
		Text:    t.Text,
	}
	return storeModel, nil
}

func (s *TemplateService) transformTemplateUpdateToStoreModel(t *models.TemplateUpdate) (*models.TemplateEntity, error) {
	uuidValue, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, fmt.Errorf("updating template: couldn't parse id: %s to UUID. Error: %w", t.ID, err)
	}
	storeModel := &models.TemplateEntity{
		ID:      uuidValue,
		Subject: t.Subject,
		Text:    t.Text,
	}
	return storeModel, nil
}
