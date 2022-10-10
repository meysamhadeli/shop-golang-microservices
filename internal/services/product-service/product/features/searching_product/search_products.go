package searching_product

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
