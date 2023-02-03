package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/consumers/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/consumers/handlers"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(
	ctx context.Context,
	jaegerTracer trace.Tracer,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config) error {

	createProductConsumer := rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, handlers.HandleConsumeCreateProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(ctx, events.ProductCreated{})
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
