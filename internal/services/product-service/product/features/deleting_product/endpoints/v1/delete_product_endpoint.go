package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/products")
	group.DELETE("/:id", deleteProduct(infra), middleware.ValidateBearerToken())
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
func deleteProduct(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v12.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			infra.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := v1.NewDeleteProduct(request.ProductID)

		if err := infra.Validator.StructCtx(ctx, command); err != nil {
			infra.Log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			infra.Log.Warn("DeleteProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
