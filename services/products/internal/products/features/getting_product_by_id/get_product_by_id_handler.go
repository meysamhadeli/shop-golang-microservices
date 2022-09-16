package getting_product_by_id

import (
	"context"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	customErrors "github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails/custome_error"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/opentracing/opentracing-go"
)

type GetProductByIdHandler struct {
	log    logger.ILogger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductByIdHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*dtos.GetProductByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		notFoundErr := customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("product with id %s not found", query.ProductID))
		return nil, notFoundErr
	}

	productDto, err := mapper.Map[*dtos.ProductResponseDto](product)
	if err != nil {
		return nil, err
	}

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
