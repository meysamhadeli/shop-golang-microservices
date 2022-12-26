package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/products")
	group.GET("/:id", getProductByID(infra), middleware.ValidateBearerToken())
}

// GetProductByID
// @Tags        Products
// @Summary     Get product
// @Description Get product by id
// @Accept      json
// @Produce     json
// @Param       id  path     string true "Product ID"
// @Success     200 {object} v1.GetProductByIdResponseDto
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [get]
func getProductByID(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v1.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			infra.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := v12.NewGetProductById(request.ProductId)

		if err := infra.Validator.StructCtx(ctx, query); err != nil {
			infra.Log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		queryResult, err := mediatr.Send[*v12.GetProductById, *v1.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			infra.Log.Warn("GetProductById", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
