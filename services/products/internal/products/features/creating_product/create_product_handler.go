package creating_product

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
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

	var commandValue, err = open_telemetry.ObjToString(command)
	if err != nil {
		return nil, err
	}
	_, span := c.jaegerTracer.Start(ctx, reflect.TypeOf(c).Elem().Name())
	span.SetAttributes(attribute.Key("productId").String(command.ProductID.String()))
	span.SetAttributes(attribute.Key("command").String(commandValue))
	defer span.End()

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

	err = c.rabbitmqPublisher.PublishMessage(evt)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
