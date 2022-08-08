package dtos

import "github.com/meysamhadeli/shop-golang-microservices/pkg/utils"

type GetProductsRequestDto struct {
	*utils.ListQuery
}
