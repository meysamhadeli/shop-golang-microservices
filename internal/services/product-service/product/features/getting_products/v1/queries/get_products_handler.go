package queries

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	dtos1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/v1/dtos"
)

type GetProductsHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository data.ProductRepository
	ctx               context.Context
	grpcClient        grpc.GrpcClient
}

func NewGetProductsHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository data.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) *GetProductsHandler {
	return &GetProductsHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher, grpcClient: grpcClient}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtos1.GetProductsResponseDto, error) {

	products, err := c.productRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)

	if err != nil {
		return nil, err
	}
	return &dtos1.GetProductsResponseDto{Products: listResultDto}, nil
}
