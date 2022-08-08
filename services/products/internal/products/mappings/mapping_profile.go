package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dto"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	return nil
}
