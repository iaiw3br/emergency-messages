package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"projects/emergency-messages/internal/errorx"
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
		s.log.Error("deleting template: uuid is empty")
		return errorx.ErrValidation
	}
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		s.log.Errorf("deleting template: invalid format uuid: %s. Error: %v", id, err)
		return errorx.ErrValidation
	}

	if err = s.templateStore.Delete(ctx, uuidValue); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.log.Errorf("deleting template: couldn't find template with id: %s", id)
			return errorx.ErrNotFound
		default:
			s.log.Errorf("deleting template: error in the server: %v", err)
			return errorx.ErrInternal
		}
	}
	return nil
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
