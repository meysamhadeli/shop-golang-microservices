package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	deleting_product2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product"
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
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete]
func (ep *deleteProductEndpoint) deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &deleting_product2.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := deleting_product2.NewDeleteProduct(request.ProductID)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*deleting_product2.DeleteProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			ep.Configuration.Log.Warn("DeleteProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
