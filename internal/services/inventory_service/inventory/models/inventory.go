package models

import (
	"time"
)

type Inventory struct {
	Id           int64     `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	ProductItems []ProductItem
}
