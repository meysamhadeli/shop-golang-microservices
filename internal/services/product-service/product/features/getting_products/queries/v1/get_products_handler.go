package queries_v1

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	dtos1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
)

type GetProductsHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo data.ProductRepository
}

func NewGetProductsHandler(log logger.ILogger, cfg *config.Config, pgRepo data.ProductRepository) *GetProductsHandler {
	return &GetProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtos1.GetProductsResponseDto, error) {

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)

	if err != nil {
		return nil, err
	}
	return &dtos1.GetProductsResponseDto{Products: listResultDto}, nil
}
