package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
	"gorm.io/gorm"
)

func migrateProducts(gorm *gorm.DB) error {

	// or we could use gorm.Migrate()
	err := gorm.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
