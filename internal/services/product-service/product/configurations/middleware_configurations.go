package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config_options"
	echo_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	echo_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	otel_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/middlewares"
	"strings"
)

func ConfigMiddlewares(e *echo_server.EchoServer, config *config_options.Config) {

	e.Echo.HideBanner = false

	e.Echo.Use(middleware.Logger())

	e.Echo.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Echo.Use(otel_middleware.EchoTracerMiddleware(config.Jaeger.ServiceName))

	e.Echo.Use(echo_middleware.CorrelationIdMiddleware)
	e.Echo.Use(middleware.RequestID())
	e.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Echo.Use(middleware.BodyLimit(constants.BodyLimit))
}
