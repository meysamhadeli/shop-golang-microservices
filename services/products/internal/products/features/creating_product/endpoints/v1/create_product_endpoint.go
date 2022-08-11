package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/dtos"
)

type createProductEndpoint struct {
	*config.ProductEndpointBase[config.InfrastructureConfiguration]
}

func NewCreteProductEndpoint(endpointBase *config.ProductEndpointBase[config.InfrastructureConfiguration]) *createProductEndpoint {
	return &createProductEndpoint{endpointBase}
}

func (ep *createProductEndpoint) MapRoute() {
	ep.ProductsGroup.POST("", ep.createProduct())
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data"
// @Success 201 {object} dtos.CreateProductResponseDto
// @Router /api/v1/products [post]
func (ep *createProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		if err := ep.Configuration.Validator.StructCtx(ctx, request); err != nil {
			ep.Configuration.Log.Errorf("(validate) err: {%v}", err)
			return err
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		result, err := mediatr.Send[*dtos.CreateProductResponseDto](ctx, command)

		if err != nil {
			ep.Configuration.Log.Errorf("(CreateProduct.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			return err
		}

		ep.Configuration.Log.Infof("(product created) id: {%s}", command.ProductID)
		return c.JSON(http.StatusCreated, result)
	}
}
