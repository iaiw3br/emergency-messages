package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"projects/emergency-messages/internal/services"
)

type User struct {
	userService services.UserService
	log         *slog.Logger
}

func NewUser(userService services.UserService, log *slog.Logger) *User {
	return &User{
		userService: userService,
		log:         log,
	}
}

func (u User) GetByCity(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	ctx := context.Background()
	users, err := u.userService.GetByCity(ctx, city)
	if assertError(err, w) {
		u.log.Error("getting by city", err)
		return
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		u.log.Error("cannot marshalling users")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(usersBytes)
	w.WriteHeader(http.StatusOK)
}

func (u User) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	usersCreated, err := u.userService.Upload(file)
	if assertError(err, w) {
		u.log.Error("uploading user", err)
		return
	}

	usersBytes, err := json.Marshal(usersCreated)
	if err != nil {
		u.log.Error("cannot marshalling users")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(usersBytes)
	w.WriteHeader(http.StatusCreated)
}
