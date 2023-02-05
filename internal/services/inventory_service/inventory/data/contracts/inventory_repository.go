package contracts

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
)

type InventoryRepository interface {
	AddProductItemToInventory(ctx context.Context, inventory *models.ProductItem) (*models.ProductItem, error)
}
