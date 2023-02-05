package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	contracts "github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data/contracts"
)

func ConfigProductsMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	inventoryRepository contracts.InventoryRepository, ctx context.Context) error {

	return nil
}
