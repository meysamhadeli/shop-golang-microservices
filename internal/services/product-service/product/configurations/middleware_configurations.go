package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echo_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	echo_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	otel_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/middlewares"
	"strings"
)

func ConfigMiddlewares(e *echo_server.EchoServer, jaegerCfg *open_telemetry.JaegerConfig) {

	e.Echo.HideBanner = false

	e.Echo.Use(middleware.Logger())

	e.Echo.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Echo.Use(otel_middleware.EchoTracerMiddleware(jaegerCfg.ServiceName))

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
