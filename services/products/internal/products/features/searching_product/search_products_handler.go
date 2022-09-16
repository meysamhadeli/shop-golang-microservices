package searching_product

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/opentracing/opentracing-go"
)

type SearchProductsHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewSearchProductsHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*dtos.SearchProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsHandler.Handle")
	defer span.Finish()

	products, err := c.pgRepo.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductResponseDto](products)
	if err != nil {
		return nil, err
	}

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
