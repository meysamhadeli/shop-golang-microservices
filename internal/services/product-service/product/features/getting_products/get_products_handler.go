package getting_products

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
)

type GetProductsHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductsHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductsHandler {
	return &GetProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtos2.GetProductsResponseDto, error) {

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos2.ProductResponseDto](products)

	if err != nil {
		return nil, err
	}
	return &dtos2.GetProductsResponseDto{Products: listResultDto}, nil
}
