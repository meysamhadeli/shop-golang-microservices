package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/consumers/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/consumers/handlers"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/shared/delivery"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(
	ctx context.Context,
	jaegerTracer trace.Tracer,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	inventoryRepository contracts.InventoryRepository,
	cfg *config.Config) error {

	inventoryDeliveryBase := delivery.InventoryDeliveryBase{
		Log:                 log,
		Cfg:                 cfg,
		JaegerTracer:        jaegerTracer,
		ConnRabbitmq:        connRabbitmq,
		InventoryRepository: inventoryRepository,
		Ctx:                 ctx,
	}

	createProductConsumer := rabbitmq.NewConsumer[*delivery.InventoryDeliveryBase](cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, handlers.HandleConsumeCreateProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(ctx, events.ProductCreated{}, &inventoryDeliveryBase)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
