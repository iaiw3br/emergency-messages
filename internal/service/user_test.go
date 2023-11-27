package service

import (
	"bytes"
	"github.com/emergency-messages/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getUsersFromCSV(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []models.User
		wantErr bool
	}{
		{
			name:    "when have valid data then no error",
			args:    "FirstName;SecondName;MobilePhone;Email\nDavid;Smith;+79216783322;david_smith@gmail.com\nLiza;Meta;89993218833;gorilla66@gmail.com",
			wantErr: false,
			want: []models.User{
				{
					FirstName:   "David",
					LastName:    "Smith",
					MobilePhone: "+79216783322",
					Email:       "david_smith@gmail.com",
				},
				{
					FirstName:   "Liza",
					LastName:    "Meta",
					MobilePhone: "89993218833",
					Email:       "gorilla66@gmail.com",
				},
			},
		},
		{
			name:    "when have less cells then error",
			args:    "FirstName;SecondName;MobilePhone\nDavid;Smith;+79216783322\nLiza;Meta;89993218833",
			wantErr: true,
			want:    nil,
		},
		{
			name:    "when have empty first name and last name then no error and no data",
			args:    "FirstName;SecondName;MobilePhone;Email\n;;+79216783322;david_smith@gmail.com\n;;89993218833;gorilla66@gmail.com",
			wantErr: false,
			want:    []models.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := bytes.NewBuffer([]byte(tt.args))
			got, err := getUsersFromCSV(data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
