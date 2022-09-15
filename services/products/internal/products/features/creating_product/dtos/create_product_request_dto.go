package dtos

type CreateProductRequestDto struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
}
