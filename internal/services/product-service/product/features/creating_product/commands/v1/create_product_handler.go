package v1

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
)

type CreateProductHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository data.ProductRepository
	ctx               context.Context
	grpcClient        grpc.GrpcClient
}

func NewCreateProductHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository data.ProductRepository, ctx context.Context, grpcClient grpc.GrpcClient) *CreateProductHandler {
	return &CreateProductHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher, grpcClient: grpcClient}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*v1.CreateProductResponseDto, error) {

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*v12.ProductCreated](createdProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(ctx, evt)
	if err != nil {
		return nil, err
	}

	response := &v1.CreateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
