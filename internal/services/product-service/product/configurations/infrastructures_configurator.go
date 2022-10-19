package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	rabbitmq2 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	events3 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
	events2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type infrastructureConfigurator struct {
	Log  logger.ILogger
	Cfg  *config.Config
	Echo *echo.Echo
}

func NewInfrastructureConfigurator(log logger.ILogger, cfg *config.Config, echo *echo.Echo) *infrastructureConfigurator {
	return &infrastructureConfigurator{Cfg: cfg, Echo: echo, Log: log}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, chan error, func()) {

	infrastructure := &shared.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, Log: ic.Log, Validator: validator.New()}

	cleanups := []func(){}

	gorm, err := ic.configGorm()
	if err != nil {
		return err, nil, nil
	}
	infrastructure.Gorm = gorm

	tp, err := ic.configOpenTelemetry()
	if err != nil {
		return err, nil, nil
	}
	infrastructure.JaegerTracer = tp.Tracer(ic.Cfg.Jaeger.TracerName)

	cleanups = append(cleanups, func() {
		err = tp.Shutdown(ctx)
		if err != nil {
			ic.Log.Fatal(err)
		}
	})

	ic.Log.Infof("%s is running", config.GetMicroserviceName(ic.Cfg.ServiceName))

	conn, err, rabbitMqCleanup := rabbitmq2.NewRabbitMQConn(ic.Cfg.Rabbitmq)
	if err != nil {
		return err, nil, nil
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	infrastructure.RabbitmqPublisher = rabbitmq2.NewPublisher(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer)

	createProductConsumer := rabbitmq2.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq2.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq2.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeDeleteProduct)

	foreverChanConsumers := make(chan error)

	go func() {
		err, createProductConsumerCleanup := createProductConsumer.ConsumeMessage(ctx, v1.ProductCreated{})
		cleanups = append(cleanups, createProductConsumerCleanup)
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		var err, updateProductConsumerCleanup = updateProductConsumer.ConsumeMessage(ctx, events2.ProductUpdated{})
		cleanups = append(cleanups, updateProductConsumerCleanup)

		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		var err, deleteProductConsumerCleanup = deleteProductConsumer.ConsumeMessage(ctx, events3.ProductDeleted{})
		cleanups = append(cleanups, deleteProductConsumerCleanup)

		if err != nil {
			ic.Log.Error(err)
		}
	}()

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	ic.configSwagger()

	ic.configMiddlewares(ic.Cfg.Jaeger)
	if err != nil {
		return err, nil, nil
	}

	pc := NewProductsModuleConfigurator(infrastructure)

	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil, nil
	}

	ic.Echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", config.GetMicroserviceName(ic.Cfg.ServiceName)))
	})

	return nil, foreverChanConsumers, func() {
		for _, c := range cleanups {
			c()
		}
	}
}
