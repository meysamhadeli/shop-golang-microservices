package v1

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
)

type GetProductsResponseDto struct {
	Products *utils.ListResult[*dtos.ProductDto]
}
