package contracts

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
	uuid "github.com/satori/go.uuid"
)

type InventoryRepository interface {
	AddProductItemToInventory(ctx context.Context, inventory *models.ProductItem) (*models.ProductItem, error)
	GetProductInInventories(ctx context.Context, productId uuid.UUID) (*models.ProductItem, error)
}
