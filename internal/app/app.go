package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	grpcapp "projects/emergency-messages/internal/app/grpc"
	"projects/emergency-messages/internal/config"
	"projects/emergency-messages/internal/consumers"
	"projects/emergency-messages/internal/controllers"
	v2 "projects/emergency-messages/internal/controllers/grpc"
	client "projects/emergency-messages/internal/databases/client/postgres"
	mdlware "projects/emergency-messages/internal/middlewares"
	"projects/emergency-messages/internal/models"
	"projects/emergency-messages/internal/providers"
	"projects/emergency-messages/internal/providers/email/mail_gun"
	"projects/emergency-messages/internal/providers/sms/twil"
	"projects/emergency-messages/internal/queue"
	"projects/emergency-messages/internal/router"
	"projects/emergency-messages/internal/senders"
	"projects/emergency-messages/internal/services"
	"projects/emergency-messages/internal/stores/postgres"
	"projects/emergency-messages/internal/workers"
	"strconv"
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

const (
	envLocal = "local"
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
		logging  = setupLogger()
		url      = os.Getenv("DATABASE_URL")
		grpcPort = os.Getenv("GRPC_PORT")
	)

	db := client.Connect(url)

	port, err := strconv.Atoi(grpcPort)
	if err != nil {
		return fmt.Errorf("cannot trasform string grpc port to integer, %w", err)
	}

	grpcApp := grpcapp.New(logging, port)
	go grpcApp.Run()

	registerEntities(db, grpcApp.GRPCServer, logging, r)

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

func setupLogger() *slog.Logger {
	env := os.Getenv("ENV")
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}
	return log
}

func registerEntities(db *bun.DB, grpcServer *grpc.Server, l *slog.Logger, r *chi.Mux) {
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
	userController := controllers.NewUser(userService, l)

	templateStore := postgres.NewTemplate(db)
	templateService := services.NewTemplate(templateStore, l)
	templateController := controllers.NewTemplate(&templateService, l)
	v2.Register(grpcServer, &templateService, l)

	brokerAddr := os.Getenv("KAFKA_BROKER")

	producer, err := queue.NewProducer(brokerAddr, l)
	if err != nil {
		log.Fatal(err)
	}

	messageStore := postgres.NewMessage(db)
	messageService := services.NewMessage(producer, templateStore, l)
	messageController := controllers.NewMessage(messageService, l)

	sender := senders.New(messageStore, userStore, l)
	messageConsumer := consumers.New(sender, messageStore, l)
	consumer, err := queue.NewConsumer(brokerAddr, messageConsumer, l)
	if err != nil {
		log.Fatal(err)
	}
	go consumer.Read()

	routers := router.New(r, messageController, userController, templateController)
	routers.Load()

	suppliers := getSuppliers(l)

	workerSendMessage := workers.NewSendMessage(messageStore, producer, suppliers, l)
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	_, err = c.AddFunc("@every 30s", workerSendMessage.Send)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}

func getSuppliers(l *slog.Logger) *providers.SendManager {
	supplier := providers.New()

	mailg := mail_gun.NewEmailMailgClient(l)
	supplier.AddProvider(mailg, models.ContactTypeEmail)

	twilSMS := twil.NewMobileTwilClient(l)
	supplier.AddProvider(twilSMS, models.ContactTypeSMS)

	return supplier
}
