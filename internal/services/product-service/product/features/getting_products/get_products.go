package getting_products

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}
