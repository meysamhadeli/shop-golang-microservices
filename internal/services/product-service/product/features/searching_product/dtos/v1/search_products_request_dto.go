package v1

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `json:"listQuery"`
}
