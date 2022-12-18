package integration

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	rabbitmq2 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"testing"
)

type TestInfrastructureConfigurator struct {
	Log logger.ILogger
	Cfg *config.Config
	t   *testing.T
}

func NewTestInfrastructureConfigurator(t *testing.T) *TestInfrastructureConfigurator {
	return &TestInfrastructureConfigurator{t: t}
}

func (ic *TestInfrastructureConfigurator) ConfigInfrastructures() (*contracts.InfrastructureConfiguration, *server.EchoServer, error, func()) {

	cleanups := []func(){}

	cfg, _ := config.InitConfig(constants.Test)

	log := logger.NewAppLogger(cfg.Logger)
	log.InitLogger()

	ctx, cancel := context.WithCancel(context.Background())

	infrastructure := &contracts.InfrastructureConfiguration{Cfg: cfg, Log: log, Validator: validator.New()}

	infrastructure.Context = ctx

	identityGrpcClient, err := grpc.NewGrpcClient(cfg.IdentityGrpcClient)

	if err != nil {
		cancel()
		return nil, nil, err, nil
	}
	infrastructure.GrpcClient = identityGrpcClient

	cleanups = append(cleanups, func() {
		_ = identityGrpcClient.Close()
	})

	gorm, err := gorm_postgres.NewGorm(cfg.GormPostgres)

	if err != nil {
		cancel()
		return nil, nil, err, nil
	}
	infrastructure.Gorm = gorm

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		cancel()
		return nil, nil, err, nil
	}

	infrastructure.ProductRepository = repositories.NewPostgresProductRepository(log, cfg, infrastructure.Gorm)

	tp, err := open_telemetry.TracerProvider(ctx, cfg.Jaeger, log)
	if err != nil {
		cancel()
		log.Fatal(err)
		return nil, nil, err, nil
	}

	infrastructure.JaegerTracer = tp.Tracer(cfg.Jaeger.TracerName)

	echoServer := server.NewEchoServer(log, cfg.Echo)
	infrastructure.Echo = echoServer.Echo
	configurations.ConfigMiddlewares(echoServer.Echo, cfg.Jaeger)

	go func() {
		if err := echoServer.RunHttpServer(ctx, nil); err != nil {
			log.Errorf("(s.RunHttpServer) err: %v", err)
		}
	}()

	conn, err, rabbitMqCleanup := rabbitmq2.NewRabbitMQTestContainers().Start(ctx, ic.t)
	if err != nil {
		cancel()
		return nil, nil, err, nil
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	infrastructure.RabbitmqPublisher = rabbitmq.NewPublisher(cfg.Rabbitmq, conn, log, infrastructure.JaegerTracer)

	if err != nil {
		cancel()
		return nil, nil, err, nil
	}

	pc := configurations.NewProductsModuleConfigurator(infrastructure, nil)

	err = configurations.ConfigureProductsModule(pc)
	if err != nil {
		cancel()
		return nil, nil, err, nil
	}

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	return infrastructure, echoServer, nil, func() {
		cancel()
		for _, c := range cleanups {
			c()
		}
	}
}
