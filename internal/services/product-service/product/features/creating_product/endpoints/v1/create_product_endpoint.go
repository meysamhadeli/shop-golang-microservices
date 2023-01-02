package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	commands_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/products")
	group.POST("", createProduct(infra), middleware.ValidateBearerToken())
}

// CreateProduct
// @Tags        Products
// @Summary     Create product
// @Description Create new product item
// @Accept      json
// @Produce     json
// @Param       CreateProductRequestDto body     v1.CreateProductRequestDto true "Product data"
// @Success     201                     {object} v1.CreateProductResponseDto
// @Security ApiKeyAuth
// @Router      /api/v1/products [post]
func createProduct(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &v1.CreateProductRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[createProductEndpoint_handler.Bind] error in the binding request")
			infra.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands_v1.NewCreateProduct(request.Name, request.Description, request.Price)

		if err := infra.Validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[createProductEndpoint_handler.StructCtx] command validation failed")
			infra.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*commands_v1.CreateProduct, *v1.CreateProductResponseDto](ctx, command)

		if err != nil {
			infra.Log.Errorf("(CreateProduct.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		infra.Log.Infof("(product created) id: {%s}", command.ProductID)
		return c.JSON(http.StatusCreated, result)
	}
}
