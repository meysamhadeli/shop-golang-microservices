package dtos

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `json:"listQuery"`
}
