package commands

import (
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewDeleteProduct(productID uuid.UUID) *DeleteProduct {
	return &DeleteProduct{ProductID: productID}
}
