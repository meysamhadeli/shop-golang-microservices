package v1

import (
	"context"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/pkg/errors"
)

type GetProductByIdHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewGetProductByIdHandler(infra *contracts.InfrastructureConfiguration) *GetProductByIdHandler {
	return &GetProductByIdHandler{infra: infra}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*v1.GetProductByIdResponseDto, error) {

	product, err := q.infra.ProductRepository.GetProductById(ctx, query.ProductID)

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
