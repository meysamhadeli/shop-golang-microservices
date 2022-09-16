package getting_products

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/opentracing/opentracing-go"
)

type GetProductsHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductsHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductsHandler {
	return &GetProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtos.GetProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductsHandler.Handle")
	defer span.Finish()

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductResponseDto](products)

	if err != nil {
		return nil, err
	}
	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
