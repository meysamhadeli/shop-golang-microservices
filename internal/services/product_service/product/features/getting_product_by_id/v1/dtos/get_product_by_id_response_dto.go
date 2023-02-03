package dtos

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/dtos"
)

type GetProductByIdResponseDto struct {
	Product *dtos.ProductDto `json:"product"`
}
