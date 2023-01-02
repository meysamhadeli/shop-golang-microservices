package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"time"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoServer struct {
	Log  logger.ILogger
	Cfg  *config.EchoConfig
	Echo *echo.Echo
}

func NewEchoServer(log logger.ILogger, cfg *config.EchoConfig) *EchoServer {
	e := echo.New()
	return &EchoServer{Log: log, Cfg: cfg, Echo: e}
}

func (s *EchoServer) RunHttpServer(ctx context.Context, configEcho ...func(echoServer *echo.Echo)) error {
	s.Echo.Server.ReadTimeout = ReadTimeout
	s.Echo.Server.WriteTimeout = WriteTimeout
	s.Echo.Server.MaxHeaderBytes = MaxHeaderBytes

	if len(configEcho) > 0 {
		httpFunc := configEcho[0]
		if httpFunc != nil {
			httpFunc(s.Echo)
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Log.Infof("shutting down Http PORT: {%s}", s.Cfg.Port)
				err := s.Echo.Shutdown(ctx)
				if err != nil {
					s.Log.Errorf("(Shutdown) err: {%v}", err)
					return
				}
				s.Log.Info("server exited properly")
				return
			}
		}
	}()

	err := s.Echo.Start(s.Cfg.Port)

	return err
}

func (s *EchoServer) ApplyVersioningFromHeader() {
	s.Echo.Pre(apiVersion)
}

func (s *EchoServer) RegisterGroupFunc(groupName string, builder func(g *echo.Group)) *EchoServer {
	builder(s.Echo.Group(groupName))

	return s
}

// APIVersion Header Based Versioning
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("/%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}
