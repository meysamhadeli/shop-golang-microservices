package dtos

import uuid "github.com/satori/go.uuid"

type ProductItemDto struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"productId"`
	Count       string    `json:"Count"`
	InventoryId uuid.UUID
}
