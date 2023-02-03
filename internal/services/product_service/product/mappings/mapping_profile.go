package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/events"
	events2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/updating_product/v1/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/models"
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

	err = mapper.CreateMap[*models.Product, *events2.ProductUpdated]()
	if err != nil {
		return err
	}
	return nil
}
