package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/commands"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/dtos"
	commands3 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/v1/commands"
	dtos4 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/v1/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/v1/queries"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/v1/dtos"
	queries2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/v1/queries"
	search_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/v1/dtos"
	queries3 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/v1/queries"
	commands2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/v1/commands"
	dtos3 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/v1/dtos"
)

func ConfigProductsMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository data.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*commands.CreateProduct, *dtos2.CreateProductResponseDto](commands.NewCreateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries2.GetProducts, *dtos.GetProductsResponseDto](queries2.NewGetProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries3.SearchProducts, *search_dtos.SearchProductsResponseDto](queries3.NewSearchProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*commands2.UpdateProduct, *dtos3.UpdateProductResponseDto](commands2.NewUpdateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*commands3.DeleteProduct, *mediatr.Unit](commands3.NewDeleteProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*queries.GetProductById, *dtos4.GetProductByIdResponseDto](queries.NewGetProductByIdHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	return nil
}
