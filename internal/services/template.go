package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	"time"
)

type TemplateService struct {
	templateStore TemplateStore
	log           *slog.Logger
}

type TemplateStore interface {
	Create(ctx context.Context, t *models.TemplateEntity) error
	Update(ctx context.Context, t *models.TemplateEntity) error
	Delete(ctx context.Context, id uuid.UUID, now time.Time) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateEntity, error)
}

func NewTemplate(templateStore TemplateStore, log *slog.Logger) TemplateService {
	return TemplateService{
		templateStore: templateStore,
		log:           log,
	}
}

func (s *TemplateService) Create(ctx context.Context, template *models.TemplateCreate) error {
	if err := template.Validate(); err != nil {
		s.log.With(slog.Any("template", template)).
			Error("validating template", err)
		return errorx.ErrValidation
	}

	storeModel, err := s.transformTemplateCreateToStoreModel(template)
	if err != nil {
		s.log.With(slog.Any("template", template)).
			Error("transforming template to store model", err)
		return errorx.ErrValidation
	}

	if err = s.templateStore.Create(ctx, storeModel); err != nil {
		s.log.With(slog.Any("template", storeModel)).
			Error("creating template", err)
		return errorx.ErrInternal
	}

	return nil
}

func (s *TemplateService) Delete(ctx context.Context, id string) error {
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("parsing uuid", id, err)
		return errorx.ErrValidation
	}

	now := time.Now()

	if err = s.templateStore.Delete(ctx, uuidValue, now); err != nil {
		s.log.With(slog.Any("templateID", id)).
			Error("deleting template", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errorx.ErrNotFound
		default:
			return errorx.ErrInternal
		}
	}
	return nil
}

func (s *TemplateService) Update(ctx context.Context, template *models.TemplateUpdate) error {
	if err := template.Validate(); err != nil {
		s.log.With(slog.Any("template", template)).
			Error("validating template", err)
		return errorx.ErrValidation
	}
	storeModel, err := s.transformTemplateUpdateToStoreModel(template)
	if err != nil {
		s.log.With(slog.Any("template", template)).
			Error("transforming template to store model", err)
		return errorx.ErrValidation
	}

	if err = s.templateStore.Update(ctx, storeModel); err != nil {
		s.log.With(slog.Any("templateID", storeModel.ID)).
			Error("updating template", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errorx.ErrNotFound
		default:
			return errorx.ErrInternal
		}
	}
	return nil
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
