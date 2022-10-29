package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echo_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	otel_middleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares(otelCfg *open_telemetry.Config) {

	ic.Echo.HideBanner = false

	ic.Echo.Use(middleware.Logger())

	ic.Echo.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	ic.Echo.Use(otel_middleware.EchoTracerMiddleware(otelCfg.ServiceName))

	ic.Echo.Use(echo_middleware.CorrelationIdMiddleware)
	ic.Echo.Use(middleware.RequestID())
	ic.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.Echo.Use(middleware.BodyLimit(constants.BodyLimit))
}
