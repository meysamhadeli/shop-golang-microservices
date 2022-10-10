package dtos

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

type GetProductsResponseDto struct {
	Products *utils.ListResult[*ProductResponseDto]
}
