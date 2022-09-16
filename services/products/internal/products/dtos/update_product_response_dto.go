package dtos

import uuid "github.com/satori/go.uuid"

type UpdateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId"`
}
