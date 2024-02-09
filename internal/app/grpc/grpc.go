package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

type App struct {
	GRPCServer *grpc.Server
	port       int
	log        *slog.Logger
}

func New(log *slog.Logger, port int) *App {
	return &App{
		GRPCServer: grpc.NewServer(),
		port:       port,
		log:        log,
	}
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatal(err)
	}

	a.log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err = a.GRPCServer.Serve(l); err != nil {
		log.Fatal(err)
	}

}

func (a *App) Stop() {
	a.GRPCServer.GracefulStop()
}
