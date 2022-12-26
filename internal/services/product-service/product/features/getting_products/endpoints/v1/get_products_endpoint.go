package endpoints_v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	queries_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/products")
	group.GET("", getAllProducts(infra), middleware.ValidateBearerToken())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param GetProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Security ApiKeyAuth
// @Router /api/v1/products [get]
func getAllProducts(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			infra.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := queries_v1.NewGetProducts(request.ListQuery)

		queryResult, err := mediatr.Send[*queries_v1.GetProducts, *dtos.GetProductsResponseDto](ctx, query)

		if err != nil {
			infra.Log.Warnf("GetProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
