package dtos

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtos.ProductDto]
}
