package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type deleteProductEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewDeleteProductEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *deleteProductEndpoint {
	return &deleteProductEndpoint{productEndpointBase}
}

func (ep *deleteProductEndpoint) MapRoute() {
	ep.ProductsGroup.DELETE("/:id", ep.deleteProduct())
}

// DeleteProduct
// @Tags        Products
// @Summary     Delete product
// @Description Delete existing product
// @Accept      json
// @Produce     json
// @Success     204
// @Param       id path string true "Product ID"
// @Router      /api/v1/products/{id} [delete]
func (ep *deleteProductEndpoint) deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v12.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := v1.NewDeleteProduct(request.ProductID)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			ep.Configuration.Log.Warn("DeleteProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
