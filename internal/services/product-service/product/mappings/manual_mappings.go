package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dto"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
)

func ProductToProductResponseDto(product *models.Product) *dto.ProductDto {
	return &dto.ProductDto{
		ProductId:   product.ProductId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
