package configurations

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	endpointsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/features/registering_user/v1/endpoints"
)

func ConfigEndpoints(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {

	endpointsv1.MapRoute(validator, log, echo, ctx)
}
