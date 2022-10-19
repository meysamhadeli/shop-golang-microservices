package v1

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

type GetProductsRequestDto struct {
	*utils.ListQuery
}
