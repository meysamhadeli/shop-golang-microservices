package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails/custome_error"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/shared"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product"
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

		request := &dtos.CreateProductRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[createProductEndpoint_handler.Bind] error in the binding request")
			ep.Configuration.Log.Errorf(fmt.Sprintf("[createProductEndpoint_handler.Bind] err: %v", badRequestErr))
			return badRequestErr
		}

		if err := ep.Configuration.Validator.StructCtx(ctx, request); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[createProductEndpoint_handler.StructCtx] command validation failed")
			ep.Configuration.Log.Errorf(fmt.Sprintf("[createProductEndpoint_handler.StructCtx] err: {%v}", validationErr))
			return validationErr
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		result, err := mediatr.Send[*creating_product.CreateProduct, *dtos.CreateProductResponseDto](ctx, command)

		if err != nil {
			ep.Configuration.Log.Errorf("(CreateProduct.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			return err
		}

		ep.Configuration.Log.Infof("(product created) id: {%s}", command.ProductID)
		return c.JSON(http.StatusCreated, result)
	}
}
