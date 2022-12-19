package integration

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
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
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

type IntegrationTestFixture struct {
	suite.Suite
	Log                 logger.ILogger
	Cfg                 *config.Config
	ProductRepository   data.ProductRepository
	IdentityGrpcClient  grpc.GrpcClient
	Ctx                 context.Context
	JaegerTracer        trace.Tracer
	RabbitMqPublisher   rabbitmq.IPublisher
	RabbitMqConsumer    rabbitmq.IConsumer
	RabbitMqConn        *amqp.Connection
	EchoServer          *server.EchoServer
	ProductEndpointBase *contracts.ProductEndpointBase[contracts.InfrastructureConfiguration]
}

func NewIntegrationTestFixture(t *testing.T) *IntegrationTestFixture {

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
		require.FailNow(t, err.Error())
	}
	infrastructure.GrpcClient = identityGrpcClient

	cleanups = append(cleanups, func() {
		_ = identityGrpcClient.Close()
	})

	gorm, err := gorm_postgres.NewGorm(cfg.GormPostgres)

	if err != nil {
		cancel()
		require.FailNow(t, err.Error())
	}
	infrastructure.Gorm = gorm

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		cancel()
		require.FailNow(t, err.Error())
	}

	infrastructure.ProductRepository = repositories.NewPostgresProductRepository(log, cfg, infrastructure.Gorm)

	tp, err := open_telemetry.TracerProvider(ctx, cfg.Jaeger, log)
	if err != nil {
		cancel()
		log.Fatal(err)
		require.FailNow(t, err.Error())
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

	conn, err, rabbitMqCleanup := rabbitmq2.NewRabbitMQTestContainers().Start(ctx, t)
	if err != nil {
		cancel()
		require.FailNow(t, err.Error())
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	infrastructure.RabbitmqPublisher = rabbitmq.NewPublisher(cfg.Rabbitmq, conn, log, infrastructure.JaegerTracer)

	if err != nil {
		cancel()
		require.FailNow(t, err.Error())
	}

	pc := configurations.NewProductsModuleConfigurator(infrastructure, nil)

	err = configurations.ConfigureProductsModule(pc)
	if err != nil {
		cancel()
		require.FailNow(t, err.Error())
	}

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	if err != nil {
		require.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		// with Cancel() we send signal to done() channel to stop  grpc, http and workers gracefully
		//https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		//https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		mediatr.ClearRequestRegistrations()
		cancel()
		for _, c := range cleanups {
			c()
		}
	})

	integration := &IntegrationTestFixture{
		Log:                 infrastructure.Log,
		Cfg:                 infrastructure.Cfg,
		Ctx:                 infrastructure.Context,
		ProductRepository:   infrastructure.ProductRepository,
		IdentityGrpcClient:  infrastructure.GrpcClient,
		JaegerTracer:        infrastructure.JaegerTracer,
		RabbitMqPublisher:   infrastructure.RabbitmqPublisher,
		RabbitMqConn:        infrastructure.ConnRabbitmq,
		ProductEndpointBase: infrastructure.ProductEndpointBase,
		EchoServer:          echoServer,
	}

	return integration
}
