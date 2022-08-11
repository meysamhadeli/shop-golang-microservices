package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/shared"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id/dtos"
)

type getProductByIdEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewGetProductByIdEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{productEndpointBase}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/:id", ep.getProductByID())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) getProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		query := getting_product_by_id.NewGetProductById(request.ProductId)

		if err := ep.Configuration.Validator.StructCtx(ctx, query); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return err
		}

		queryResult, err := mediatr.Send[*getting_product_by_id.GetProductById, *dtos.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warn("GetProductById", err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
