package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projects/emergency-messages/internal/config"
	client "projects/emergency-messages/internal/databases/client/postgres"
	"projects/emergency-messages/internal/handlers"
	"projects/emergency-messages/internal/logging"
	mdlware "projects/emergency-messages/internal/middlewares"
	mailg "projects/emergency-messages/internal/providers/email/mail_gun"
	"projects/emergency-messages/internal/services"
	"projects/emergency-messages/internal/stores/postgres"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uptrace/bun"
)

var (
	shutdownTimeout = 5 * time.Second
	contextTimeout  = 60 * time.Second
)

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
		url     = os.Getenv("DATABASE_URL")
	)

	db := client.Connect(url)

	registerEntities(db, logging, r)

	// if err := migration.RunMigrate(url); err != nil {
	// 	log.Fatalf("running migration: %v", err)
	// }

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func registerEntities(db *bun.DB, l logging.Logger, r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mdlware.LimitRequests)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(contextTimeout))

	userStore := postgres.NewUserStore(db)
	userService := services.NewUserService(userStore, l)
	userHandler := handlers.NewUser(userService, l)
	userHandler.Register(r)

	templateStore := postgres.NewTemplate(db)
	templateService := services.NewTemplate(templateStore, l)
	templateHandler := handlers.NewTemplate(templateService, l)
	templateHandler.Register(r)

	mailg := mailg.New(l)

	messageStore := postgres.NewMessage(db)
	messageService := services.NewMessage(messageStore, templateStore, userStore, mailg, l)
	messageHandler := handlers.NewMessage(messageService, l)
	messageHandler.Register(r)
}
