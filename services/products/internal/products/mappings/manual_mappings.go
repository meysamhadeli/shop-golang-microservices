package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
)

func ProductToProductResponseDto(product *models.Product) *dtos.ProductResponseDto {
	return &dtos.ProductResponseDto{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
