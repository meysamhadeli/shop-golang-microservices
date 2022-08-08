package infrastructure

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	product_constants "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/constants"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/web/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares(metrics *CatalogsServiceMetrics) {

	ic.echo.HideBanner = false

	ic.echo.HTTPErrorHandler = middlewares.ProblemHandler

	middlewareManager := middlewares.NewMiddlewareManager(ic.log, ic.cfg, getHttpMetricsCb(metrics))

	ic.echo.Use(middlewareManager.RequestLoggerMiddleware)
	ic.echo.Use(middlewareManager.RequestMetricsMiddleware)

	ic.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         product_constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	ic.echo.Use(middleware.RequestID())
	ic.echo.Use(middleware.Logger())
	ic.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: product_constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.echo.Use(middleware.BodyLimit(product_constants.BodyLimit))
}

func getHttpMetricsCb(metrics *CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
