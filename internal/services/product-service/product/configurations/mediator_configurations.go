package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	v17 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	create_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v13 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	v15 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	v14 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	queries_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	queries_v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/queries/v1"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

func ConfigProductsMediator(infra *contracts.InfrastructureConfiguration) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v17.CreateProduct, *create_dtos.CreateProductResponseDto](v17.NewCreateProductHandler(infra))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries_v1.GetProducts, *dtos.GetProductsResponseDto](queries_v1.NewGetProductsHandler(infra))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries_v12.SearchProducts, *search_dtos.SearchProductsResponseDto](queries_v12.NewSearchProductsHandler(infra))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v1.UpdateProduct, *v12.UpdateProductResponseDto](v1.NewUpdateProductHandler(infra))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v13.DeleteProduct, *mediatr.Unit](v13.NewDeleteProductHandler(infra))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v14.GetProductById, *v15.GetProductByIdResponseDto](v14.NewGetProductByIdHandler(infra))
	if err != nil {
		return err
	}

	return nil
}
