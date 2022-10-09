package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echo2 "github.com/meysamhadeli/shop-golang-microservices/pkg/http/echo/config"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"time"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoServer struct {
	Log    logger.ILogger
	Cfg    *echo2.EchoConfig
	Echo   *echo.Echo
	DoneCh chan struct{}
}

func NewEchoServer(log logger.ILogger, cfg *echo2.EchoConfig) *EchoServer {
	return &EchoServer{Log: log, Cfg: cfg, Echo: echo.New()}
}

func (s *EchoServer) RunHttpServer(configEcho func(echoServer *echo.Echo)) error {
	s.Echo.Server.ReadTimeout = ReadTimeout
	s.Echo.Server.WriteTimeout = WriteTimeout
	s.Echo.Server.MaxHeaderBytes = MaxHeaderBytes

	if configEcho != nil {
		configEcho(s.Echo)
	}

	return s.Echo.Start(s.Cfg.Port)
}

func (s *EchoServer) ApplyVersioningFromHeader() {
	s.Echo.Pre(apiVersion)
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
