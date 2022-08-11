package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/gorm_postgres"
	"gorm.io/gorm"
)

func (ic *infrastructureConfigurator) configGorm() (*gorm.DB, error) {
	gorm, err := gorm_postgres.NewGorm(ic.Cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	return gorm, nil
}
