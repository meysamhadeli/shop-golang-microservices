package dtos

import uuid "github.com/satori/go.uuid"

type RegisterUserResponseDto struct {
	UserId    uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
