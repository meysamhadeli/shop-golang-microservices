package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"gorm.io/gorm"
)

func migrateProducts(gorm *gorm.DB) error {

	// or we could use gorm.Migrate()
	err := gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
