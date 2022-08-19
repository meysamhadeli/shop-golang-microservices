package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/interceptors"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/consume"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/shared"
	"google.golang.org/grpc"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type infrastructureConfigurator struct {
	Log        logger.ILogger
	Cfg        *config.Config
	Echo       *echo.Echo
	GrpcServer *grpc.Server
}

func NewInfrastructureConfigurator(log logger.ILogger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *infrastructureConfigurator {
	return &infrastructureConfigurator{Cfg: cfg, Echo: echo, GrpcServer: grpcServer, Log: log}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, func()) {

	infrastructure := &shared.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, GrpcServer: ic.GrpcServer, Log: ic.Log, Validator: validator.New()}

	infrastructure.Im = interceptors.NewInterceptorManager(ic.Log)

	cleanups := []func(){}

	gorm, err := ic.configGorm()
	if err != nil {
		return err, nil
	}
	infrastructure.Gorm = gorm

	conn, err, rabbitMqCleanup := rabbitmq.NewRabbitMQConn(ic.Cfg.Rabbitmq)
	if err != nil {
		return err, nil
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	rabbitMqPublisher := rabbitmq.NewPublisher(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log)
	infrastructure.RabbitmqPublisher = rabbitMqPublisher

	createProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, consume.HandleConsumeCreateProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(events.ProductCreated{})
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	ic.configSwagger()

	ic.configMiddlewares()

	if err != nil {
		return err, nil
	}

	pc := NewProductsModuleConfigurator(infrastructure)
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil
	}

	ic.Echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", config.GetMicroserviceName(ic.Cfg)))
	})

	return nil, func() {
		for _, c := range cleanups {
			defer c()
		}
	}
}
