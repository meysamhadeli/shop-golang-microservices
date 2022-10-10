package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"github.com/pkg/errors"
	"net/http"
)

type updateProductEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewUpdateProductEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *updateProductEndpoint {
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

		request := &dtos.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[updateProductEndpoint_handler.Bind] error in the binding request")
			ep.Configuration.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := updating_product.NewUpdateProduct(request.ProductID, request.Name, request.Description, request.Price)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[updateProductEndpoint_handler.StructCtx] command validation failed")
			ep.Configuration.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*updating_product.UpdateProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			ep.Configuration.Log.Warnf("UpdateProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		ep.Configuration.Log.Infof("(product updated) id: {%s}", request.ProductID)

		return c.NoContent(http.StatusNoContent)
	}
}
