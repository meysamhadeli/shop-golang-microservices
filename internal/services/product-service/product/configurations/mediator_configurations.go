package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	v17 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	create_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v16 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	get_by_id_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	v15 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/queries/v1"
	get_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	v14 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/queries/v1"
	v13 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v17.CreateProduct, *create_dtos.CreateProductResponseDto](v17.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher, c.JaegerTracer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v1.GetProducts, *get_dtos.GetProductsResponseDto](v1.NewGetProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v14.SearchProducts, *search_dtos.SearchProductsResponseDto](v14.NewSearchProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v13.UpdateProduct, *v12.UpdateProductResponseDto](v13.NewUpdateProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v16.DeleteProduct, *mediatr.Unit](v16.NewDeleteProductHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v15.GetProductById, *get_by_id_dtos.GetProductByIdResponseDto](v15.NewGetProductByIdHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	return nil
}
