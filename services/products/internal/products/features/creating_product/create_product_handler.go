package creating_product

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"time"

	kafkaClient "github.com/meysamhadeli/shop-golang-microservices/pkg/kafka"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/tracing"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts/grpc/kafka_messages"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/proto"
)

type CreateProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	repository    contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, kafkaProducer: kafkaProducer}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtos.CreateProductResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	defer span.Finish()

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

	evt := &kafka_messages.ProductCreated{Product: mappings.ProductToGrpcMessage(createdProduct)}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	err = c.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductID}
	bytes, _ := json.Marshal(response)

	span.LogFields(log.String("CreateProductResponseDto", string(bytes)))

	return response, nil
}
