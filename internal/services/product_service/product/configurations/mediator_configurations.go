package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/data/contracts"
	creatingproductv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/commands"
	creatingproductv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/dtos"
	deletingproductv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/deleting_product/v1/commands"
	gettingproductbyidv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_product_by_id/v1/dtos"
	gettingproductbyidv1queries "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_product_by_id/v1/queries"
	gettingproductsv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_products/v1/dtos"
	gettingproductsv1queries "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/getting_products/v1/queries"
	searchingproductv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/searching_product/v1/dtos"
	searchingproductv1queries "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/searching_product/v1/queries"
	updatingproductv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/updating_product/v1/commands"
	updatingproductv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/updating_product/v1/dtos"
)

func ConfigProductsMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingproductv1commands.CreateProduct, *creatingproductv1dtos.CreateProductResponseDto](creatingproductv1commands.NewCreateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingproductsv1queries.GetProducts, *gettingproductsv1dtos.GetProductsResponseDto](gettingproductsv1queries.NewGetProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searchingproductv1queries.SearchProducts, *searchingproductv1dtos.SearchProductsResponseDto](searchingproductv1queries.NewSearchProductsHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updatingproductv1commands.UpdateProduct, *updatingproductv1dtos.UpdateProductResponseDto](updatingproductv1commands.NewUpdateProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deletingproductv1commands.DeleteProduct, *mediatr.Unit](deletingproductv1commands.NewDeleteProductHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingproductbyidv1queries.GetProductById, *gettingproductbyidv1dtos.GetProductByIdResponseDto](gettingproductbyidv1queries.NewGetProductByIdHandler(log, rabbitmqPublisher, productRepository, ctx, grpcClient))
	if err != nil {
		return err
	}

	return nil
}
