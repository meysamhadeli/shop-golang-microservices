package v1

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type RegisterUser struct {
	UserId    uuid.UUID `validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	UserName  string    `json:"userName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=4"`
	CreatedAt time.Time `validate:"required"`
}

func NewRegisterUser(firstName string, lastName string, userName string, email string, password string) *RegisterUser {
	return &RegisterUser{UserId: uuid.NewV4(), FirstName: firstName, LastName: lastName, UserName: userName, Email: email, Password: password, CreatedAt: time.Now()}
}
