package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	events2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dtos.ProductResponseDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events2.ProductCreated]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events2.ProductUpdated]()
	if err != nil {
		return err
	}
	return nil
}
