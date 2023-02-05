package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Inventory, *dtos.InventoryDto]()
	if err != nil {
		return err
	}
	return nil
}
