package controller

import (
	"fmt"
	"github.com/emergency-messages/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const users = "/users"

type UserController struct {
	UserService service.UserService
}

func (u UserController) Register(r *chi.Mux) {
	r.Get(users, u.GetAll)
	r.Post(fmt.Sprintf("%s/upload", users), u.Upload)
}

func NewUser(userService service.UserService) UserController {
	return UserController{
		UserService: userService,
	}
}

func (u UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all users"))
}

func (u UserController) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = u.UserService.Upload(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
