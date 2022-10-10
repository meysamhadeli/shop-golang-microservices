package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	creating_product2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product"
	deleting_product2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product"
	getting_product_by_id2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id"
	getting_products2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products"
	searching_product2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product"
	updating_product2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creating_product2.CreateProduct, *dtos2.CreateProductResponseDto](creating_product2.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher, c.JaegerTracer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getting_products2.GetProducts, *dtos2.GetProductsResponseDto](getting_products2.NewGetProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searching_product2.SearchProducts, *dtos2.SearchProductsResponseDto](searching_product2.NewSearchProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updating_product2.UpdateProduct, *dtos2.UpdateProductResponseDto](updating_product2.NewUpdateProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deleting_product2.DeleteProduct, *mediatr.Unit](deleting_product2.NewDeleteProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getting_product_by_id2.GetProductById, *dtos2.GetProductByIdResponseDto](getting_product_by_id2.NewGetProductByIdHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	return nil
}
