package models

import "testing"

func TestTemplateUpdate_Validate(t1 *testing.T) {
	type fields struct {
		ID      uint64
		Subject string
		Text    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "when id is empty then error",
			fields: fields{
				ID:      0,
				Subject: "1",
				Text:    "2",
			},
			wantErr: true,
		},
		{
			name: "when subject is  empty then error",
			fields: fields{
				ID:   1,
				Text: "2",
			},
			wantErr: true,
		},
		{
			name: "when text is empty then error",
			fields: fields{
				ID:      0,
				Subject: "1",
			},
			wantErr: true,
		},
		{
			name: "when all data are not empty then no error",
			fields: fields{
				ID:      1,
				Subject: "1",
				Text:    "2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TemplateUpdate{
				ID:      tt.fields.ID,
				Subject: tt.fields.Subject,
				Text:    tt.fields.Text,
			}
			if err := t.Validate(); (err != nil) != tt.wantErr {
				t1.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
