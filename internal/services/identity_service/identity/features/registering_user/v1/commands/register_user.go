package commands

import (
	"time"
)

type RegisterUser struct {
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	UserName  string    `json:"userName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=4"`
	CreatedAt time.Time `validate:"required"`
}

func NewRegisterUser(firstName string, lastName string, userName string, email string, password string) *RegisterUser {
	return &RegisterUser{FirstName: firstName, LastName: lastName, UserName: userName, Email: email, Password: password, CreatedAt: time.Now()}
}
