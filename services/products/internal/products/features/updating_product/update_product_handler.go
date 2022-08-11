package updating_product

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/http_errors"
	kafkaClient "github.com/meysamhadeli/shop-golang-microservices/pkg/kafka"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts/grpc/kafka_messages"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type UpdateProductHandler struct {
	log           logger.ILogger
	cfg           *config.Config
	pgRepo        contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewUpdateProductHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*mediatr.Unit, error) {

	_, err := c.pgRepo.GetProductById(ctx, command.ProductID)

	if err != nil {
		return nil, http_errors.NewNotFoundError(fmt.Sprintf("product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt := &kafka_messages.ProductUpdated{Product: mappings.ProductToGrpcMessage(updatedProduct)}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}

	message := kafka.Message{
		Topic: c.cfg.KafkaTopics.ProductUpdated.TopicName,
		Value: msgBytes,
		Time:  time.Now(),
	}

	return &mediatr.Unit{}, c.kafkaProducer.PublishMessage(ctx, message)
}
