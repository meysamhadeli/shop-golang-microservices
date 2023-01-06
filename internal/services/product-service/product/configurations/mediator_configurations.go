package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	v17 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	create_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v13 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/commands/v1"
	v15 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	v14 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	queries_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/dtos/v1"
	queries_v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/queries/v1"
	update_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/commands/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
)

func ConfigProductsMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository data.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v17.CreateProduct, *create_dtos.CreateProductResponseDto](v17.NewCreateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries_v1.GetProducts, *v1.GetProductsResponseDto](queries_v1.NewGetProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries_v12.SearchProducts, *search_dtos.SearchProductsResponseDto](queries_v12.NewSearchProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*update_v1.UpdateProduct, *v12.UpdateProductResponseDto](update_v1.NewUpdateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v13.DeleteProduct, *mediatr.Unit](v13.NewDeleteProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*v14.GetProductById, *v15.GetProductByIdResponseDto](v14.NewGetProductByIdHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	return nil
}
