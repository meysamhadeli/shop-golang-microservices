package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	events "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dtos.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events.ProductCreated]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *v1.ProductUpdated]()
	if err != nil {
		return err
	}
	return nil
}
