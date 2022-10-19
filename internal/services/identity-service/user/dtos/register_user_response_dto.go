package dtos

type RegisterUserResponseDto struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
