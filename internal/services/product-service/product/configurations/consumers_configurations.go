package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

func ConfigConsumers(infra *contracts.InfrastructureConfiguration) error {

	createProductConsumer := rabbitmq.NewConsumer(infra.Cfg, infra.ConnRabbitmq, infra.Log, infra.JaegerTracer, consumers2.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(infra.Cfg, infra.ConnRabbitmq, infra.Log, infra.JaegerTracer, consumers2.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(infra.Cfg, infra.ConnRabbitmq, infra.Log, infra.JaegerTracer, consumers2.HandleConsumeDeleteProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(infra.Context, v1.ProductCreated{})
		if err != nil {
			infra.Log.Error(err)
		}
	}()

	go func() {
		err := updateProductConsumer.ConsumeMessage(infra.Context, v12.ProductUpdated{})
		if err != nil {
			infra.Log.Error(err)
		}
	}()

	go func() {
		err := deleteProductConsumer.ConsumeMessage(infra.Context, events.ProductDeleted{})
		if err != nil {
			infra.Log.Error(err)
		}
	}()

	return nil
}
