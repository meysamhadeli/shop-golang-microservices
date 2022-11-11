package v1

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
)

type SearchProductsHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo data.ProductRepository
}

func NewSearchProductsHandler(log logger.ILogger, cfg *config.Config, pgRepo data.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*search_dtos.SearchProductsResponseDto, error) {

	products, err := c.pgRepo.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &search_dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
