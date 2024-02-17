package models

import (
	"github.com/google/uuid"
	"testing"
)

func TestMessageRequest_Validate(t *testing.T) {
	type fields struct {
		TemplateID uuid.UUID
		City       string
		Strength   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				TemplateID: uuid.New(),
				City:       "city",
				Strength:   "strength",
			},
			wantErr: false,
		},
		{
			name: "invalid template id",
			fields: fields{
				TemplateID: uuid.Nil,
				City:       "city",
				Strength:   "strength",
			},
			wantErr: true,
		},
		{
			name: "invalid city",
			fields: fields{
				TemplateID: uuid.New(),
				City:       "",
				Strength:   "strength",
			},
			wantErr: true,
		},
		{
			name: "invalid strength",
			fields: fields{
				TemplateID: uuid.New(),
				City:       "city",
				Strength:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageRequest{
				TemplateID: tt.fields.TemplateID,
				City:       tt.fields.City,
				Strength:   tt.fields.Strength,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
