package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"gorm.io/gorm"
)

func (ic *infrastructureConfigurator) configGorm() (*gorm.DB, error) {
	gorm, err := gorm_postgres.NewGorm(ic.Cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	err = migrateProducts(gorm)
	if err != nil {
		return nil, err
	}

	return gorm, nil
}

func migrateProducts(gorm *gorm.DB) error {

	// or we could use gorm.Migrate()
	err := gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
