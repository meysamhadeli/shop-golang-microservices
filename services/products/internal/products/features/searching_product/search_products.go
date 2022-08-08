package searching_product

import "github.com/meysamhadeli/shop-golang-microservices/pkg/utils"

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
