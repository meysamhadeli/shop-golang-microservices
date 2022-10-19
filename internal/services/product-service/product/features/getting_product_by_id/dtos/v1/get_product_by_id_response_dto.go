package v1

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
)

type GetProductByIdResponseDto struct {
	Product *dtos.ProductDto `json:"product"`
}
