package queries_v1

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

type SearchProductsHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewSearchProductsHandler(infra *contracts.InfrastructureConfiguration) *SearchProductsHandler {
	return &SearchProductsHandler{infra: infra}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*search_dtos.SearchProductsResponseDto, error) {

	products, err := c.infra.ProductRepository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &search_dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
