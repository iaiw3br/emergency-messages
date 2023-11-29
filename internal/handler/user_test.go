package handler

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserController_Upload(t *testing.T) {
	u := UserController{}
	server := httptest.NewServer(http.HandlerFunc(u.GetAll))

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := "all users"
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, want, string(b))
}

func TestUserController_Upload1(t *testing.T) {
	t.Skip()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("err not expected while open a mock db, %v", err)
	}
	defer db.Close()

	mock.
		ExpectQuery("INSERT INTO users (first_name, last_name, mobile_phone, email) VALUES ($1, $2, $3, $4)").
		WithArgs("boris", "smith", "+7921462398323", "john.doe@example.com").
		WillReturnError(pgx.ErrNoRows)

	f, err := os.Open("file.csv")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/upload", f)
	// r.Header.Set("Content-Type", "multipart/form-data")
	// r.Header.Set("Content-Type", "multipart/form-data")
	// r.Header.Set("Content-Disposition", "form-data; name=\"file\"; filename=\"file.csv\"")
	// r.Header.Set("Content-Type", "Content-Type: text/csv")
	ctx := context.WithValue(r.Context(), "DB", db)
	r = r.WithContext(ctx)

	// request Content-Type isn't multipart/form-data
	// no multipart boundary param in Content-Type
	u := UserController{}
	u.Upload(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("not all expectations were met: %v", err)
	}
}
