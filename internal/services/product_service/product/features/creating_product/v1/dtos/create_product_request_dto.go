package dtos

type CreateProductRequestDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Count       int32   `json:"count"`
	InventoryId int64   `json:"inventoryId"`
}
