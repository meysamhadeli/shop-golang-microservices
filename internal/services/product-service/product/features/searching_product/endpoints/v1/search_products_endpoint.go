package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	query_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type searchProductsEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewSearchProductsEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *searchProductsEndpoint {
	return &searchProductsEndpoint{productEndpointBase}
}

func (ep *searchProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/search", ep.searchProducts())
}

// SearchProducts
// @Tags        Products
// @Summary     Search products
// @Description Search products
// @Accept      json
// @Produce     json
// @Param       searchProductsRequestDto query    v1.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success     200                      {object} v1.SearchProductsResponseDto[dto.ProductDto]
// @Router      /api/v1/products/search [get]
func (ep *searchProductsEndpoint) searchProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &v1.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := &query_v1.SearchProducts{SearchText: request.SearchText, ListQuery: request.ListQuery}

		if err := ep.Configuration.Validator.StructCtx(ctx, query); err != nil {
			ep.Configuration.Log.Errorf("(validate) err: {%v}", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		queryResult, err := mediatr.Send[*query_v1.SearchProducts, *v1.SearchProductsResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warn("SearchProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
