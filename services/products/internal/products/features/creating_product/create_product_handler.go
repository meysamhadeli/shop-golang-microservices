package creating_product

import (
	"context"
	"encoding/json"
	kafkaClient "github.com/meysamhadeli/shop-golang-microservices/pkg/kafka"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
)

type CreateProductHandler struct {
	log               logger.ILogger
	cfg               *config.Config
	repository        contracts.ProductRepository
	kafkaProducer     kafkaClient.Producer
	rabbitmqPublisher rabbitmq.IPublisher
}

func NewCreateProductHandler(log logger.ILogger, cfg *config.Config, repository contracts.ProductRepository, kafkaProducer kafkaClient.Producer,
	rabbitmqPublisher rabbitmq.IPublisher) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, kafkaProducer: kafkaProducer, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtos.CreateProductResponseDto, error) {

	product := &models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt := &events.ProductCreated{ProductId: createdProduct.ProductID}

	err = c.rabbitmqPublisher.PublishMessage(evt)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductID}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
