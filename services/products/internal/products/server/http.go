package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/consts"

	"time"
)

func (s *Server) RunHttpServer(configEcho func(echoServer *echo.Echo)) error {
	s.Echo.Server.ReadTimeout = consts.ReadTimeout
	s.Echo.Server.WriteTimeout = consts.WriteTimeout
	s.Echo.Server.MaxHeaderBytes = consts.MaxHeaderBytes

	if configEcho != nil {
		configEcho(s.Echo)
	}

	return s.Echo.Start(s.Cfg.Http.Port)
}

func (s *Server) WaitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.DoneCh <- struct{}{}
	}()
}

func (s *Server) ApplyVersioningFromHeader() {
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
