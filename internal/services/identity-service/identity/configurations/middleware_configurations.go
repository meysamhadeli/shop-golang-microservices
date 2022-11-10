package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echo_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	otel_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/middlewares"
	"strings"
)

func configMiddlewares(e *echo.Echo, otelCfg *open_telemetry.Config) {

	e.HideBanner = false

	e.Use(middleware.Logger())
	e.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Use(otel_middleware.EchoTracerMiddleware(otelCfg.ServiceName))

	e.Use(echo_middleware.CorrelationIdMiddleware)
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit(constants.BodyLimit))
}
