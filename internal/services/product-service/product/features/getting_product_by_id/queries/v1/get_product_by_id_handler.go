package v1

import (
	"context"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/dtos/v1"
	"github.com/pkg/errors"
)

type GetProductByIdHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository data.ProductRepository
	ctx               context.Context
	grpcClient        grpc.GrpcClient
}

func NewGetProductByIdHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository data.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher, grpcClient: grpcClient}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*v1.GetProductByIdResponseDto, error) {

	product, err := q.productRepository.GetProductById(ctx, query.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", query.ProductID))
		return nil, notFoundErr
	}

	productDto, err := mapper.Map[*dtos.ProductDto](product)
	if err != nil {
		return nil, err
	}

	return &v1.GetProductByIdResponseDto{Product: productDto}, nil
}
