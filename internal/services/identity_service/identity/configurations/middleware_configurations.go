package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echomiddleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel"
	otelmiddleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/middlewares"
	"strings"
)

func ConfigMiddlewares(e *echo.Echo, jaegerCfg *otel.JaegerConfig) {

	e.HideBanner = false

	e.Use(middleware.Logger())
	e.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Use(otelmiddleware.EchoTracerMiddleware(jaegerCfg.ServiceName))

	e.Use(echomiddleware.CorrelationIdMiddleware)
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit(constants.BodyLimit))
}
