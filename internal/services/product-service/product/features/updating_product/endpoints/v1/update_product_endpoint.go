package v1

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	commands_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	"github.com/pkg/errors"
	"net/http"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.PUT("/:id", updateProduct(validator, log, ctx), middleware.ValidateBearerToken())
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
func updateProduct(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := &v1.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[updateProductEndpoint_handler.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands_v1.NewUpdateProduct(request.ProductId, request.Name, request.Description, request.Price)

		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[updateProductEndpoint_handler.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*commands_v1.UpdateProduct, *v1.UpdateProductResponseDto](ctx, command)

		if err != nil {
			log.Warnf("UpdateProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("(product updated) id: {%s}", request.ProductId)

		return c.NoContent(http.StatusNoContent)
	}
}
