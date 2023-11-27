package app

import (
	"context"
	"github.com/emergency-messages/internal/config"
	"github.com/emergency-messages/internal/controller"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/service"
	"github.com/emergency-messages/internal/store"
	"github.com/emergency-messages/pkg/client/postgres"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// logging +
// config +
// server +
// database

func Run() {
	_ = logging.New()

	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// database
	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	r := chi.NewRouter()

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore)
	users := controller.NewUser(userService)
	users.Register(r)
	// port := os.Getenv("PORT")
	http.ListenAndServe(":8080", r)
}
