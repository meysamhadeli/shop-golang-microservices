package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product"
	creating_products_dtos "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/deleting_product"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id"
	getting_by_id_dtos "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_products"
	getting_products_dtos "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/searching_product"
	searching_products_dtos "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/searching_product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/updating_product"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creating_product.CreateProduct, *creating_products_dtos.CreateProductResponseDto](creating_product.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer, c.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getting_products.GetProducts, *getting_products_dtos.GetProductsResponseDto](getting_products.NewGetProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searching_product.SearchProducts, *searching_products_dtos.SearchProductsResponseDto](searching_product.NewSearchProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updating_product.UpdateProduct, *mediatr.Unit](updating_product.NewUpdateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deleting_product.DeleteProduct, *mediatr.Unit](deleting_product.NewDeleteProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getting_product_by_id.GetProductById, *getting_by_id_dtos.GetProductByIdResponseDto](getting_product_by_id.NewGetProductByIdHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	return nil
}
