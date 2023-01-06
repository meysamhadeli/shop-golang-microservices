package v1

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/dtos/v1"
	"net/http"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.DELETE("/:id", deleteProduct(validator, log, ctx), middleware.ValidateBearerToken())
}

// DeleteProduct
// @Tags        Products
// @Summary     Delete product
// @Description Delete existing product
// @Accept      json
// @Produce     json
// @Success     204
// @Param       id path string true "Product ID"
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [delete]
func deleteProduct(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := &v12.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := v1.NewDeleteProduct(request.ProductID)

		if err := validator.StructCtx(ctx, command); err != nil {
			log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			log.Warn("DeleteProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
