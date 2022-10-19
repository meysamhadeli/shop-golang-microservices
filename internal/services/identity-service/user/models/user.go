package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// User model
type User struct {
	UserId    uuid.UUID `json:"userId" gorm:"primaryKey"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
