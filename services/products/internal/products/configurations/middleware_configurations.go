package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/constants"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares(otelCfg *open_telemetry.Config) {

	ic.Echo.HideBanner = false

	ic.Echo.Use(middleware.Logger())

	ic.Echo.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	ic.Echo.Use(middlewares.EchoTracerMiddleware(otelCfg.ServiceName))

	ic.Echo.Use(middlewares.CorrelationIdMiddleware)
	ic.Echo.Use(middleware.RequestID())
	ic.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.Echo.Use(middleware.BodyLimit(constants.BodyLimit))
}
