package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	events2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/v1/events"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/v1/events"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(
	ctx context.Context,
	jaegerTracer trace.Tracer,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config) error {

	createProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers2.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers2.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers2.HandleConsumeDeleteProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(ctx, events2.ProductCreated{})
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		err := updateProductConsumer.ConsumeMessage(ctx, v12.ProductUpdated{})
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		err := deleteProductConsumer.ConsumeMessage(ctx, events.ProductDeleted{})
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
