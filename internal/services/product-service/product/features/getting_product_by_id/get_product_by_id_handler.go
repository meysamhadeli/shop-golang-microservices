package getting_product_by_id

import (
	"context"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/pkg/errors"
)

type GetProductByIdHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductByIdHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*dtos2.GetProductByIdResponseDto, error) {

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", query.ProductID))
		return nil, notFoundErr
	}

	productDto, err := mapper.Map[*dtos2.ProductResponseDto](product)
	if err != nil {
		return nil, err
	}

	return &dtos2.GetProductByIdResponseDto{Product: productDto}, nil
}
