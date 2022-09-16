package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dtos.ProductResponseDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events.ProductCreated]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events.ProductUpdated]()
	if err != nil {
		return err
	}
	return nil
}
