package integration

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/context_provider"
	ech_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	logger "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	gorm2 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/gorm"
	rabbitmq2 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
	"os"
	"testing"
)

type IntegrationTestFixture struct {
	suite.Suite
	Log               logger.ILogger
	Cfg               *config.Config
	RabbitmqPublisher rabbitmq.IPublisher
	RabbitmqConsumer  *rabbitmq.Consumer
	ConnRabbitmq      *amqp.Connection
	HttpClient        *resty.Client
	JaegerTracer      trace.Tracer
	Gorm              *gorm.DB
	Echo              *echo.Echo
	GrpcClient        grpc.GrpcClient
	ProductRepository data.ProductRepository
	Context           context.Context
}

func NewIntegrationTestFixture(t *testing.T, option fx.Option) *IntegrationTestFixture {

	err := os.Setenv("APP_ENV", constants.Test)

	if err != nil {
		require.FailNow(t, err.Error())
	}

	integrationTestFixture := &IntegrationTestFixture{}

	app := fxtest.New(t,
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				context_provider.NewContext,
				ech_server.NewEchoServer,
				grpc.NewGrpcClient,
				gorm_postgres.NewGorm,
				open_telemetry.TracerProvider,
				http_client.NewHttpClient,
				repositories.NewPostgresProductRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				configurations.InitialInfrastructures,
			),
			fx.Invoke(func(infrastructure *contracts.InfrastructureConfiguration) {

				//// get gorm-postgres from test-container
				gormDb, err := gorm2.NewGormTestContainers().Start(infrastructure.Context, t)
				if err != nil {
					require.FailNow(t, err.Error())
				}

				// get rabbitmq from test-container
				connRabbitMq, err := rabbitmq2.NewRabbitMQTestContainers().Start(infrastructure.Context, t)
				if err != nil {
					require.FailNow(t, err.Error())
				}

				integrationTestFixture.Gorm = gormDb
				integrationTestFixture.ConnRabbitmq = connRabbitMq

				integrationTestFixture.Log = infrastructure.Log
				integrationTestFixture.Cfg = infrastructure.Cfg
				integrationTestFixture.RabbitmqPublisher = infrastructure.RabbitmqPublisher
				integrationTestFixture.HttpClient = infrastructure.HttpClient
				integrationTestFixture.JaegerTracer = infrastructure.JaegerTracer
				integrationTestFixture.Echo = infrastructure.Echo
				integrationTestFixture.GrpcClient = infrastructure.GrpcClient
				integrationTestFixture.ProductRepository = infrastructure.ProductRepository
				integrationTestFixture.Context = infrastructure.Context
			}),
			fx.Invoke(func(gorm *gorm.DB) error {
				return gorm_postgres.Migrate(gorm, &models.Product{})
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigProductsMediator),
			option,
		),
	)

	defer app.RequireStart().RequireStop()
	require.NoError(t, app.Err())

	configurations.ConfigMiddlewares(&ech_server.EchoServer{Echo: integrationTestFixture.Echo}, integrationTestFixture.Cfg.Jaeger)

	return integrationTestFixture
}
