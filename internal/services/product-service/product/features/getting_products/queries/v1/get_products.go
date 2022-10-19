package v1

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
)

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

func NewGetProducts(query *utils.ListQuery) *GetProducts {
	return &GetProducts{ListQuery: query}
}
