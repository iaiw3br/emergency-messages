package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const users = "/users"

type UserController struct {
	userService service.UserService
	log         logging.Logger
}

func (u UserController) Register(r *chi.Mux) {
	r.Get(fmt.Sprintf("%s/city/:city", users), u.GetByCity)
	r.Post(fmt.Sprintf("%s/upload", users), u.Upload)
}

func NewUser(userService service.UserService, log logging.Logger) UserController {
	return UserController{
		userService: userService,
		log:         log,
	}
}

func (u UserController) GetByCity(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	ctx := context.Background()
	users, err := u.userService.GetByCity(ctx, city)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		u.log.Error("cannot marshal users")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(usersBytes)
	w.WriteHeader(http.StatusOK)
}

func (u UserController) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	usersCreated, err := u.userService.Upload(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usersBytes, err := json.Marshal(usersCreated)
	if err != nil {
		u.log.Error("cannot marshal users")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(usersBytes)
	w.WriteHeader(http.StatusCreated)
}
