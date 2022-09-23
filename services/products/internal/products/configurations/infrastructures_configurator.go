package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/consumers"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/shared"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
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

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, chan error, *tracesdk.TracerProvider, func()) {

	infrastructure := &shared.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, Log: ic.Log, Validator: validator.New()}

	cleanups := []func(){}

	gorm, err := ic.configGorm()
	if err != nil {
		return err, nil, nil, func() {
			for _, c := range cleanups {
				defer c()
			}
		}
	}
	infrastructure.Gorm = gorm

	tp, err := ic.configOpenTelemetry()
	infrastructure.JaegerTracer = tp.Tracer(ic.Cfg.Jaeger.TracerName)

	if err != nil {
		return err, nil, nil, func() {
			for _, c := range cleanups {
				defer c()
			}
		}
	}

	conn, err, rabbitMqCleanup := rabbitmq.NewRabbitMQConn(ic.Cfg.Rabbitmq)
	if err != nil {
		return err, nil, nil, func() {
			for _, c := range cleanups {
				defer c()
			}
		}
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	infrastructure.RabbitmqPublisher = rabbitmq.NewPublisher(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer)

	createProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers.HandleConsumeDeleteProduct)

	foreverChanConsumers := make(chan error)

	go func() {
		err, createProductConsumerCleanup := createProductConsumer.ConsumeMessage(ctx, events.ProductCreated{})
		cleanups = append(cleanups, createProductConsumerCleanup)
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		var err, updateProductConsumerCleanup = updateProductConsumer.ConsumeMessage(ctx, events.ProductUpdated{})
		cleanups = append(cleanups, updateProductConsumerCleanup)

		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		var err, deleteProductConsumerCleanup = deleteProductConsumer.ConsumeMessage(ctx, events.ProductDeleted{})
		cleanups = append(cleanups, deleteProductConsumerCleanup)

		if err != nil {
			ic.Log.Error(err)
		}
	}()

	ic.configSwagger()

	ic.configMiddlewares()

	if err != nil {
		return err, nil, nil, func() {
			for _, c := range cleanups {
				defer c()
			}
		}
	}

	pc := NewProductsModuleConfigurator(infrastructure)
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil, nil, func() {
			for _, c := range cleanups {
				defer c()
			}
		}
	}

	ic.Echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", config.GetMicroserviceName(ic.Cfg)))
	})

	return nil, foreverChanConsumers, tp, func() {
		for _, c := range cleanups {
			defer c()
		}
	}
}
