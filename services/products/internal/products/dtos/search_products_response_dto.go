package dtos

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*ProductResponseDto]
}
