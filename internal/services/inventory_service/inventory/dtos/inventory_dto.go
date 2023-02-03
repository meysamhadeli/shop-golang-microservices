package dtos

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type InventoryDto struct {
	Id           uuid.UUID        `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	ProductItems []ProductItemDto `json:"productItems"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}
