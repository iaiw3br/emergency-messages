package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mock_controllers "projects/emergency-messages/internal/controllers/mocks"
	"projects/emergency-messages/internal/errorx"
	"projects/emergency-messages/internal/models"
	api "projects/emergency-messages/protos"
	"testing"
)

func Test_template_Delete(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mock_controllers.NewMockTemplateService(controller)
	temp := template{templateService: service}

	ctx := context.Background()

	t.Run("when all queue is invalid then no error", func(t *testing.T) {
		req := &api.DeleteRequest{Id: "abc"}

		service.EXPECT().
			Delete(ctx, req.GetId()).
			Return(nil)

		resp, err := temp.Delete(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("when store returns not found error then error", func(t *testing.T) {
		req := &api.DeleteRequest{Id: "abc"}

		service.EXPECT().
			Delete(ctx, req.GetId()).
			Return(errorx.ErrNotFound)

		resp, err := temp.Delete(ctx, req)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.NotFound, s.Code())
		assert.Equal(t, "not found", s.Message())
	})

	t.Run("when store returns internal error then error", func(t *testing.T) {
		req := &api.DeleteRequest{Id: "abc"}

		service.EXPECT().
			Delete(ctx, req.GetId()).
			Return(errorx.ErrInternal)

		resp, err := temp.Delete(ctx, req)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.Internal, s.Code())
		assert.Equal(t, "template delete with error", s.Message())
	})
}

func Test_template_Create(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock_controllers.NewMockTemplateService(controller)

	ctx := context.Background()
	t.Run("when all queue is invalid then no error", func(t *testing.T) {
		tmpl := &api.CreateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateCreate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(nil)

		temp := template{templateService: service}
		resp, err := temp.Create(ctx, tmpl)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("when store returns validation error then error", func(t *testing.T) {
		tmpl := &api.CreateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateCreate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(errorx.ErrValidation)

		temp := template{templateService: service}
		resp, err := temp.Create(ctx, tmpl)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.InvalidArgument, s.Code())
		assert.Equal(t, "invalid input queue", s.Message())
	})

	t.Run("when store returns internal error then error", func(t *testing.T) {
		tmpl := &api.CreateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateCreate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Create(ctx, newTemplate).
			Return(errorx.ErrInternal)

		temp := template{templateService: service}
		resp, err := temp.Create(ctx, tmpl)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.Internal, s.Code())
		assert.Equal(t, "error while creating", s.Message())
	})
}

func Test_template_Update(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock_controllers.NewMockTemplateService(controller)

	ctx := context.Background()
	t.Run("when all queue is valid then no error", func(t *testing.T) {
		tmpl := &api.UpdateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateUpdate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(nil)

		temp := template{templateService: service}
		resp, err := temp.Update(ctx, tmpl)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("when all queue is invalid then error", func(t *testing.T) {
		tmpl := &api.UpdateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateUpdate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(errorx.ErrValidation)

		temp := template{templateService: service}
		resp, err := temp.Update(ctx, tmpl)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.InvalidArgument, s.Code())
		assert.Equal(t, "invalid input queue", s.Message())
	})

	t.Run("when all queue is not found then error", func(t *testing.T) {
		tmpl := &api.UpdateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateUpdate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(errorx.ErrNotFound)

		temp := template{templateService: service}
		resp, err := temp.Update(ctx, tmpl)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.NotFound, s.Code())
		assert.Equal(t, "not found", s.Message())
	})

	t.Run("when all queue is internal error then error", func(t *testing.T) {
		tmpl := &api.UpdateRequest{
			Text:    "321",
			Subject: "subject",
		}

		newTemplate := &models.TemplateUpdate{
			Text:    tmpl.Text,
			Subject: tmpl.Subject,
		}

		service.EXPECT().
			Update(ctx, newTemplate).
			Return(errorx.ErrInternal)

		temp := template{templateService: service}
		resp, err := temp.Update(ctx, tmpl)
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		s, ok := status.FromError(err)
		assert.Equal(t, true, ok)
		assert.Equal(t, codes.Internal, s.Code())
		assert.Equal(t, "error while updating", s.Message())
	})
}
