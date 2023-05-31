package commands

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/data/contracts"
	eventsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/deleting_product/v1/events"
)

type DeleteProductHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
	grpcClient        grpc.GrpcClient
}

func NewDeleteProductHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) *DeleteProductHandler {
	return &DeleteProductHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher, grpcClient: grpcClient}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {

	if err := c.productRepository.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	err := c.rabbitmqPublisher.PublishMessage(eventsv1.ProductDeleted{
		ProductId: command.ProductID,
	})
	if err != nil {
		return nil, err
	}

	c.log.Info("DeleteProduct successfully executed")

	return &mediatr.Unit{}, err
}
