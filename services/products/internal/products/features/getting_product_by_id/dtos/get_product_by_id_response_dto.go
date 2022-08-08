package dtos

import "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
