package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/constants"
	middlewares "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares() {

	ic.Echo.HideBanner = false

	ic.Echo.HTTPErrorHandler = middlewares.ProblemDetailsHandler

	middlewareManager := middlewares.NewMiddlewareManager(ic.Log, ic.Cfg)

	ic.Echo.Use(middlewareManager.RequestLoggerMiddleware)

	ic.Echo.Use(middleware.RequestID())
	ic.Echo.Use(middleware.Logger())
	ic.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.Echo.Use(middleware.BodyLimit(constants.BodyLimit))
}
