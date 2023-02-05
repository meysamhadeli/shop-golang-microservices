package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProductItem struct {
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	ProductId   uuid.UUID `json:"productId"`
	Count       int32     `json:"count"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	InventoryId int64
}
