package endpoints

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	echomiddleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	dtosv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_products/v1/dtos"
	queriesv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_products/v1/queries"
	"net/http"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.GET("", getAllProducts(validator, log, ctx), echomiddleware.ValidateBearerToken())
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
func getAllProducts(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &dtosv1.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := queriesv1.NewGetProducts(request.ListQuery)

		queryResult, err := mediatr.Send[*queriesv1.GetProducts, *dtosv1.GetProductsResponseDto](ctx, query)

		if err != nil {
			log.Warnf("GetProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
