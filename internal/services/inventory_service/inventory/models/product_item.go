package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ProductItem struct {
	gorm.Model
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	ProductId   uuid.UUID `json:"productId"`
	Count       string    `json:"Count"`
	InventoryId uuid.UUID
}
