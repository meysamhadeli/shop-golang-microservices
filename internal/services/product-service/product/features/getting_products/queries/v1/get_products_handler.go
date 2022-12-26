package queries_v1

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	dtos1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

type GetProductsHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewGetProductsHandler(infra *contracts.InfrastructureConfiguration) *GetProductsHandler {
	return &GetProductsHandler{infra: infra}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtos1.GetProductsResponseDto, error) {

	products, err := c.infra.ProductRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)

	if err != nil {
		return nil, err
	}
	return &dtos1.GetProductsResponseDto{Products: listResultDto}, nil
}
