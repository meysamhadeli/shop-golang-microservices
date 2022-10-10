package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"github.com/pkg/errors"
	"net/http"
)

type createProductEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewCreteProductEndpoint(endpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *createProductEndpoint {
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
// @Success 201 dtos.CreateProductResponseDto
// @Router /api/v1/products [post]
func (ep *createProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &dtos2.CreateProductRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[createProductEndpoint_handler.Bind] error in the binding request")
			ep.Configuration.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := ep.Configuration.Validator.StructCtx(ctx, request); err != nil {
			validationErr := errors.Wrap(err, "[createProductEndpoint_handler.StructCtx] command validation failed")
			ep.Configuration.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		result, err := mediatr.Send[*creating_product.CreateProduct, *dtos2.CreateProductResponseDto](ctx, command)

		if err != nil {
			ep.Configuration.Log.Errorf("(CreateProduct.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		ep.Configuration.Log.Infof("(product created) id: {%s}", command.ProductID)
		return c.JSON(http.StatusCreated, result)
	}
}
