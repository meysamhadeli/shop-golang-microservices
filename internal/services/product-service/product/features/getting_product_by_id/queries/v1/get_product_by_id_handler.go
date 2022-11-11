package v1

import (
	"context"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	"github.com/pkg/errors"
)

type GetProductByIdHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo data.ProductRepository
}

func NewGetProductByIdHandler(log logger.ILogger, cfg *config.Config, pgRepo data.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*v1.GetProductByIdResponseDto, error) {

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", query.ProductID))
		return nil, notFoundErr
	}

	productDto, err := mapper.Map[*dtos.ProductDto](product)
	if err != nil {
		return nil, err
	}

	return &v1.GetProductByIdResponseDto{Product: productDto}, nil
}
