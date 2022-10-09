package server

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/configurations"
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

	echoServer := server.NewEchoServer(s.Log, s.Cfg.Echo)

	go func() {
		if err := echoServer.RunHttpServer(nil); err != nil {
			s.Log.Errorf("(s.RunHttpServer) err: {%v}", err)
			cancel()
		}
	}()

	infrastructureConfigurator := configurations.NewInfrastructureConfigurator(s.Log, s.Cfg, echoServer.Echo)
	err, doneChanConsumers, tp, productsCleanup := infrastructureConfigurator.ConfigInfrastructures(ctx)
	if err != nil {
		return err
	}

	defer productsCleanup()

	<-ctx.Done()
	<-doneChanConsumers
	<-echoServer.DoneCh

	err = tp.Shutdown(ctx)
	if err != nil {
		s.Log.Fatal(err)
	}

	s.Log.Infof("%s is shutting down Http PORT: {%s}", config.GetMicroserviceName(s.Cfg.ServiceName), s.Cfg.Echo.Port)
	if err := echoServer.Echo.Shutdown(ctx); err != nil {
		s.Log.Warnf("(Shutdown) err: {%v}", err)
	}
	s.Log.Infof("%s server exited properly", config.GetMicroserviceName(s.Cfg.ServiceName))

	return nil
}
