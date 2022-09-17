package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/constants"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/configurations"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Log        logger.ILogger
	Cfg        *config.Config
	Echo       *echo.Echo
	DoneCh     chan struct{}
	GrpcServer *grpc.Server
}

func NewServer(log logger.ILogger, cfg *config.Config) *Server {
	return &Server{Log: log, Cfg: cfg, Echo: echo.New()}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	infrastructureConfigurator := configurations.NewInfrastructureConfigurator(s.Log, s.Cfg, s.Echo, s.GrpcServer)
	err, productsCleanup, doneChanConsumers := infrastructureConfigurator.ConfigInfrastructures(ctx)
	if err != nil {
		return err
	}
	defer productsCleanup()

	go func() {
		if err := s.RunHttpServer(nil); err != nil {
			s.Log.Errorf("(s.RunHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.Log.Infof("%s is listening on Http PORT: {%s}", config.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

	<-ctx.Done()
	s.WaitShootDown(constants.WaitShotDownDuration)

	s.Log.Infof("%s is shutting down Http PORT: {%s}", config.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)
	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Log.Warnf("(Shutdown) err: {%v}", err)
	}

	<-doneChanConsumers
	<-s.DoneCh
	s.Log.Infof("%s server exited properly", config.GetMicroserviceName(s.Cfg))

	return nil
}
