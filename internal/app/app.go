package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/internal/config"
	"github.com/emergency-messages/internal/handler"
	"github.com/emergency-messages/internal/logging"
	mdlware "github.com/emergency-messages/internal/middleware"
	"github.com/emergency-messages/internal/providers/email/mailgun"
	"github.com/emergency-messages/internal/service"
	"github.com/emergency-messages/internal/store"
	"github.com/emergency-messages/pkg/client/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	)

	db, err := postgres.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	registerEntities(db, logging, r)

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

func registerEntities(db *pgx.Conn, l logging.Logger, r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mdlware.LimitRequests)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(contextTimeout))

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore, l)
	userHandler := handler.NewUser(userService, l)
	userHandler.Register(r)

	templateStore := store.NewTemplate(db, l)
	templateService := service.NewTemplate(templateStore, l)
	templateHandler := handler.NewTemplate(templateService, l)
	templateHandler.Register(r)

	mailg := mailg.New(l)

	messageStore := store.NewMessage(db, l)
	messageService := service.NewMessage(messageStore, templateStore, userStore, mailg, l)
	messageHandler := handler.NewMessage(messageService, l)
	messageHandler.Register(r)
}
