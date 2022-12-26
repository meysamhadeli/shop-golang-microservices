package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	commands_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/products")
	group.PUT("/:id", updateProduct(infra), middleware.ValidateBearerToken())
}

// UpdateProduct
// @Tags        Products
// @Summary     Update product
// @Description Update existing product
// @Accept      json
// @Produce     json
// @Param       UpdateProductRequestDto body v1.UpdateProductRequestDto true "Product data"
// @Param       id                      path string                       true "Product ID"
// @Success     204
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [put]
func updateProduct(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v1.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[updateProductEndpoint_handler.Bind] error in the binding request")
			infra.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands_v1.NewUpdateProduct(request.ProductId, request.Name, request.Description, request.Price)

		if err := infra.Validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[updateProductEndpoint_handler.StructCtx] command validation failed")
			infra.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*commands_v1.UpdateProduct, *v1.UpdateProductResponseDto](ctx, command)

		if err != nil {
			infra.Log.Warnf("UpdateProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		infra.Log.Infof("(product updated) id: {%s}", request.ProductId)

		return c.NoContent(http.StatusNoContent)
	}
}
