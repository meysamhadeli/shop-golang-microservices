package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/

type UpdateProductRequestDto struct {
	ProductId   uuid.UUID `json:"-" param:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required"`
	Count       int32     `json:"count" validate:"required"`
	InventoryId int64     `json:"inventoryId" validate:"required"`
}
