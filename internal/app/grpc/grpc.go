package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"projects/emergency-messages/internal/logging"
)

type App struct {
	GRPCServer *grpc.Server
	port       int
	log        logging.Logger
}

func New(log logging.Logger, port int) *App {
	return &App{
		GRPCServer: grpc.NewServer(),
		port:       port,
		log:        log,
	}
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Fatal(err)
	}

	a.log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err = a.GRPCServer.Serve(l); err != nil {
		a.log.Fatal(err)
	}

}

func (a *App) Stop() {
	a.GRPCServer.GracefulStop()
}
