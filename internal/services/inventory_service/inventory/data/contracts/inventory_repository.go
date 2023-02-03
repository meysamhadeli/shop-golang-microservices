package contracts

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
)

type InventoryRepository interface {
	AddProductsToInventory(ctx context.Context, inventory *models.Inventory) (*models.Inventory, error)
}
