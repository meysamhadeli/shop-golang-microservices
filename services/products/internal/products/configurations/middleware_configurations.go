package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	product_constants "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/consts"
	middlewares2 "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares() {

	ic.Echo.HideBanner = false

	ic.Echo.HTTPErrorHandler = middlewares2.ProblemHandler

	middlewareManager := middlewares2.NewMiddlewareManager(ic.Log, ic.Cfg)

	ic.Echo.Use(middlewareManager.RequestLoggerMiddleware)
	ic.Echo.Use(middlewareManager.RequestMetricsMiddleware)

	ic.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         product_constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	ic.Echo.Use(middleware.RequestID())
	ic.Echo.Use(middleware.Logger())
	ic.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: product_constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.Echo.Use(middleware.BodyLimit(product_constants.BodyLimit))
}
