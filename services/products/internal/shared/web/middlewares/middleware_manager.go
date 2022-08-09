package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	RequestMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log       logger.ILogger
	cfg       *config.Config
	metricsCb MiddlewareMetricsCb
}

func NewMiddlewareManager(log logger.ILogger, cfg *config.Config, metricsCb MiddlewareMetricsCb) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg, metricsCb: metricsCb}
}
