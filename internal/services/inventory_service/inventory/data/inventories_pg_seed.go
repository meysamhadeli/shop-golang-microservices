package data

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func SeedData(gorm *gorm.DB) error {

	var rowsAffected = gorm.First(&models.Inventory{}).RowsAffected

	if rowsAffected == 0 {
		err := gorm.CreateInBatches(inventoriesSeeds, len(inventoriesSeeds)).Error
		if err != nil {
			return errors.Wrap(err, "error in seed inventories database")
		}
	}

	return nil
}

var inventoriesSeeds = []*models.Inventory{
	{
		Id:          1,
		Name:        "food",
		Description: "some food inventories data",
		CreatedAt:   time.Now(),
	},
	{
		Id:          2,
		Name:        "health",
		Description: "some health inventories data",
		CreatedAt:   time.Now(),
	},
}
