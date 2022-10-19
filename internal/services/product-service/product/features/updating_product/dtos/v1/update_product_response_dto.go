package v1

import uuid "github.com/satori/go.uuid"

type UpdateProductResponseDto struct {
	ProductId uuid.UUID `json:"productId"`
}
