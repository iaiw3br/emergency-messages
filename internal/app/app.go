package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/config"
	"github.com/emergency-messages/internal/controller"
	"github.com/emergency-messages/internal/logging"
	"github.com/emergency-messages/internal/service"
	"github.com/emergency-messages/internal/store"
	"github.com/emergency-messages/pkg/client/postgres"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownTimeout = 5 * time.Second

func Run() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := startServer(ctx); err != nil {
		log.Fatal(err)
	}
}

func startServer(ctx context.Context) error {
	var (
		r          = chi.NewRouter()
		listenAddr = os.Getenv("LISTEN_ADDR")
		srv        = &http.Server{
			Addr:    listenAddr,
			Handler: r,
		}
		logging = logging.New()
	)

	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore, logging)
	users := controller.NewUser(userService)
	users.Register(r)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	log.Printf("listening on %s", listenAddr)
	<-ctx.Done()

	log.Println("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown: %w", ctx.Err())
	case <-longShutdown:
		log.Println("finished")
	}

	return nil
}
