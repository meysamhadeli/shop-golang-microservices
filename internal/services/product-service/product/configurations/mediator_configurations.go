package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	v17 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	create_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v16 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	get_by_id_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	v15 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/queries/v1"
	get_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	queries_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	v14 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/queries/v1"
	v13 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

func ConfigProductsMediator(ic *contracts.InfrastructureConfiguration) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v17.CreateProduct, *create_dtos.CreateProductResponseDto](v17.NewCreateProductHandler(ic.Log, ic.Cfg, ic.ProductRepository, ic.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries_v1.GetProducts, *get_dtos.GetProductsResponseDto](queries_v1.NewGetProductsHandler(ic.Log, ic.Cfg, ic.ProductRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v14.SearchProducts, *search_dtos.SearchProductsResponseDto](v14.NewSearchProductsHandler(ic.Log, ic.Cfg, ic.ProductRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v13.UpdateProduct, *v12.UpdateProductResponseDto](v13.NewUpdateProductHandler(ic.Log, ic.Cfg, ic.ProductRepository, ic.RabbitmqPublisher, ic.GrpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v16.DeleteProduct, *mediatr.Unit](v16.NewDeleteProductHandler(ic.Log, ic.Cfg, ic.ProductRepository, ic.RabbitmqPublisher))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v15.GetProductById, *get_by_id_dtos.GetProductByIdResponseDto](v15.NewGetProductByIdHandler(ic.Log, ic.Cfg, ic.ProductRepository))
	if err != nil {
		return err
	}

	return nil
}
