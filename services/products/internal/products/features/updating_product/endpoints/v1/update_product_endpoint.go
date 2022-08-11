package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/updating_product"
)

type updateProductEndpoint struct {
	*config.ProductEndpointBase[config.InfrastructureConfiguration]
}

func NewUpdateProductEndpoint(productEndpointBase *config.ProductEndpointBase[config.InfrastructureConfiguration]) *updateProductEndpoint {
	return &updateProductEndpoint{productEndpointBase}
}

func (ep *updateProductEndpoint) MapRoute() {
	ep.ProductsGroup.PUT("/:id", ep.updateProduct())
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body updating_product.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /api/v1/products/{id} [put]
func (ep *updateProductEndpoint) updateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &updating_product.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		command := updating_product.NewUpdateProduct(request.ProductID, request.Name, request.Description, request.Price)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return err
		}

		_, err := mediatr.Send[*mediatr.Unit](ctx, command)

		if err != nil {
			ep.Configuration.Log.Warnf("UpdateProduct", err)
			return err
		}

		ep.Configuration.Log.Infof("(product updated) id: {%s}", request.ProductID)

		return c.NoContent(http.StatusNoContent)
	}
}
