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

type getProductByIdEndpoint struct {
	*contracts.ProductEndpointBase[contracts.InfrastructureConfiguration]
}

func NewGetProductByIdEndpoint(productEndpointBase *contracts.ProductEndpointBase[contracts.InfrastructureConfiguration]) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{productEndpointBase}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/:id", ep.getProductByID(), middleware.ValidateBearerToken())
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
func (ep *getProductByIdEndpoint) getProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v1.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := v12.NewGetProductById(request.ProductId)

		if err := ep.Configuration.Validator.StructCtx(ctx, query); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		queryResult, err := mediatr.Send[*v12.GetProductById, *v1.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warn("GetProductById", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
