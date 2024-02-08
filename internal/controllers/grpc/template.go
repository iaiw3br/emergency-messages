package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/logging"
	"projects/emergency-messages/internal/models"
	api "projects/emergency-messages/protos"
)

type TemplateService interface {
	Create(ctx context.Context, template *models.TemplateCreate) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, template *models.TemplateUpdate) error
}

type template struct {
	templateService TemplateService
	api.UnimplementedTemplateServer
	log logging.Logger
}

func Register(grpcServer *grpc.Server, templateService TemplateService, log logging.Logger) {
	api.RegisterTemplateServer(
		grpcServer,
		&template{
			templateService: templateService,
			log:             log,
		},
	)
}

func (t *template) Delete(ctx context.Context, req *api.DeleteRequest) (*api.EmptyResponse, error) {
	err := t.templateService.Delete(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "template delete with error")
	}
	return &api.EmptyResponse{}, nil
}

func (t *template) Create(ctx context.Context, req *api.CreateRequest) (*api.EmptyResponse, error) {
	newTemplate := &models.TemplateCreate{
		Subject: req.GetSubject(),
		Text:    req.GetText(),
	}
	if err := t.templateService.Create(ctx, newTemplate); err != nil {
		if errors.Is(err, errorx.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, "invalid input data")
		}
		return nil, status.Error(codes.Internal, "error while creating")
	}
	return &api.EmptyResponse{}, nil
}

func (t *template) Update(ctx context.Context, req *api.UpdateRequest) (*api.EmptyResponse, error) {
	updateTemplate := &models.TemplateUpdate{
		Subject: req.GetSubject(),
		Text:    req.GetText(),
	}
	if err := t.templateService.Update(ctx, updateTemplate); err != nil {
		switch {
		case errors.Is(err, errorx.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, "invalid input data")
		case errors.Is(err, errorx.ErrNotFound):
			return nil, status.Error(codes.NotFound, "not found")
		default:
			return nil, status.Error(codes.Internal, "error while updating")
		}
	}
	return &api.EmptyResponse{}, nil
}
