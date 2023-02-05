package dtos

import uuid "github.com/satori/go.uuid"

type ProductItemDto struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"productId"`
	Count       int32     `json:"count"`
	InventoryId int64     `json:"inventoryId"`
}
