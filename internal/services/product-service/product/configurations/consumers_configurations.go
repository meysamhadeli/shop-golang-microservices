package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(ctx context.Context, cfg *config.Config, conn *amqp.Connection, log logger.ILogger, tracer trace.Tracer) error {

	createProductConsumer := rabbitmq.NewConsumer(cfg, conn, log, tracer, consumers2.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(cfg, conn, log, tracer, consumers2.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(cfg, conn, log, tracer, consumers2.HandleConsumeDeleteProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(ctx, v1.ProductCreated{})
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
