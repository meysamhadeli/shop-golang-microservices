package server

import (
	"context"
	grpc "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc/server"
	echo "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/configurations"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Log logger.ILogger
	Cfg *config.Config
}

func NewServer(log logger.ILogger, cfg *config.Config) *Server {
	return &Server{Log: log, Cfg: cfg}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	echoServer := echo.NewEchoServer(s.Log, s.Cfg.Echo)
	grpcServer := grpc.NewGrpcServer(s.Log, s.Cfg.Grpc)

	go func() {
		if err := echoServer.RunHttpServer(ctx); err != nil {
			cancel()
			s.Log.Fatalf("(s.RunHttpServer) err: {%v}", err)
		}
	}()

	go func() {
		if err := grpcServer.RunGrpcServer(ctx); err != nil {
			cancel()
			s.Log.Fatalf("(s.RunGrpcServer) err: {%v}", err)
		}
	}()

	infrastructureConfigurator := configurations.NewInfrastructureConfigurator(s.Log, s.Cfg, echoServer.Echo)
	err, identitiesCleanup := infrastructureConfigurator.ConfigInfrastructures(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	defer func() {
		identitiesCleanup()
	}()

	return err
}
