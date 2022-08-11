package middlewares

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
)

type middlewareManager struct {
	log logger.ILogger
	cfg *config.Config
}

func NewMiddlewareManager(log logger.ILogger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
}
