package store

import (
	"github.com/emergency-messages/internal/models"
	"github.com/emergency-messages/internal/store/mock_store"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTemplate_Create(t *testing.T) {
	t.Run("when send all data have then no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		template := models.Template{
			Subject: "MSCH",
			Text:    "be careful",
		}
		temp := mock_store.NewMockTemplater(ctrl)
		temp.
			EXPECT().
			Create(gomock.Any(), template).
			Return(nil).
			AnyTimes()
	})
	t.Run("when send nil then error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		temp := mock_store.NewMockTemplater(ctrl)
		temp.
			EXPECT().
			Create(gomock.Any(), nil).
			Return(nil).
			AnyTimes()
	})
}

func TestTemplate_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	temp := mock_store.NewMockTemplater(ctrl)
	temp.
		EXPECT().
		Delete(gomock.Any(), gomock.Eq(10)).
		Return(nil).
		AnyTimes()
}

func TestTemplate_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	temp := mock_store.NewMockTemplater(ctrl)
	var id uint64 = 10
	temp.
		EXPECT().
		GetByID(gomock.Any(), id).
		Return(&models.Template{}, nil).
		AnyTimes()
}

func TestTemplate_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	template := models.Template{
		Subject: "MSCH",
		Text:    "be careful",
	}
	temp := mock_store.NewMockTemplater(ctrl)
	temp.
		EXPECT().
		Update(gomock.Any(), template).
		Return(nil).
		AnyTimes()
}
