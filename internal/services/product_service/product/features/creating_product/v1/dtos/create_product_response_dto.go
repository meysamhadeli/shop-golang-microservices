package dtos

import uuid "github.com/satori/go.uuid"

type CreateProductResponseDto struct {
	ProductId uuid.UUID `json:"productId"`
}
