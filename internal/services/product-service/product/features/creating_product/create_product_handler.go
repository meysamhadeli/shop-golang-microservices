package creating_product

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"go.opentelemetry.io/otel/trace"
)

type CreateProductHandler struct {
	log               logger.ILogger
	cfg               *config.Config
	repository        contracts.ProductRepository
	rabbitmqPublisher rabbitmq.IPublisher
	jaegerTracer      trace.Tracer
}

func NewCreateProductHandler(log logger.ILogger, cfg *config.Config, repository contracts.ProductRepository,
	rabbitmqPublisher rabbitmq.IPublisher, jaegerTracer trace.Tracer) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, rabbitmqPublisher: rabbitmqPublisher, jaegerTracer: jaegerTracer}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtos.CreateProductResponseDto, error) {

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*events.ProductCreated](createdProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(ctx, evt)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
