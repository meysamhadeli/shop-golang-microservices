package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Inventory struct {
	gorm.Model
	Id           uuid.UUID     `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	ProductItems []ProductItem `json:"productItems"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}
