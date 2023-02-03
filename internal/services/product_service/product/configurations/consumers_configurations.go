package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/consumers"
	creatingproducteventsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/deleting_product/v1/events"
	updatingproducteventsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/updating_product/v1/events"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(
	ctx context.Context,
	jaegerTracer trace.Tracer,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config) error {

	createProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers.HandleConsumeDeleteProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(ctx, creatingproducteventsv1.ProductCreated{})
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		err := updateProductConsumer.ConsumeMessage(ctx, updatingproducteventsv1.ProductUpdated{})
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
