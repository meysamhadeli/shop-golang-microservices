package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_products"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_products/dtos"
)

type getProductsEndpoint struct {
	*config.ProductEndpointBase[config.InfrastructureConfiguration]
}

func NewGetProductsEndpoint(productEndpointBase *config.ProductEndpointBase[config.InfrastructureConfiguration]) *getProductsEndpoint {
	return &getProductsEndpoint{productEndpointBase}
}

func (ep *getProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("", ep.getAllProducts())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, ep.Configuration.Log, err)
			return err
		}

		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		query := &getting_products.GetProducts{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*dtos.GetProductsResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warnf("GetProducts", err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
