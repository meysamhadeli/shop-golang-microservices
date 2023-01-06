package configurations

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	creating_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/endpoints/v1"
	deleting_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/endpoints/v1"
	getting_product_by_id "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/endpoints/v1"
	getting_products "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/endpoints/v1"
	searching_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/endpoints/v1"
	updating_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/endpoints/v1"
)

func ConfigEndpoints(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {

	creating_product.MapRoute(validator, log, echo, ctx)
	deleting_product.MapRoute(validator, log, echo, ctx)
	getting_product_by_id.MapRoute(validator, log, echo, ctx)
	getting_products.MapRoute(validator, log, echo, ctx)
	searching_product.MapRoute(validator, log, echo, ctx)
	updating_product.MapRoute(validator, log, echo, ctx)
}
