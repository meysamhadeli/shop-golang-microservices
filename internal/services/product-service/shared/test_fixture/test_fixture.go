package test_fixture

import (
	"context"
	"github.com/go-playground/validator"
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
	gorm_container "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/gorm"
	rabbitmq_container "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
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

type TestFixture struct {
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

func NewIntegrationTestFixture(t *testing.T, option fx.Option) *TestFixture {

	err := os.Setenv("APP_ENV", constants.Test)

	if err != nil {
		require.FailNow(t, err.Error())
	}

	integrationTestFixture := &TestFixture{}

	app := fxtest.New(t,
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				context_provider.NewContext,
				ech_server.NewEchoServer,
				gorm_postgres.NewGorm,
				grpc.NewGrpcClient,
				open_telemetry.TracerProvider,
				http_client.NewHttpClient,
				repositories.NewPostgresProductRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(func(
				rabbitmqPublisher rabbitmq.IPublisher,
				productRepository data.ProductRepository,
				ctx context.Context,
				grpcClient grpc.GrpcClient,
				echo *echo.Echo,
				log logger.ILogger,
				jaegerTracer trace.Tracer,
				httpClient *resty.Client,
				validator *validator.Validate,
				cfg *config.Config,
			) {

				// get gorm-postgres from test-container
				gormDb, err := gorm_container.NewGormTestContainers().Start(ctx, t)
				if err != nil {
					require.FailNow(t, err.Error())
				}

				// get rabbitmq from test-container
				connRabbitMq, err := rabbitmq_container.NewRabbitMQTestContainers().Start(ctx, t)
				if err != nil {
					require.FailNow(t, err.Error())
				}

				integrationTestFixture.Gorm = gormDb
				integrationTestFixture.ConnRabbitmq = connRabbitMq

				integrationTestFixture.Log = log
				integrationTestFixture.Cfg = cfg
				integrationTestFixture.RabbitmqPublisher = rabbitmqPublisher
				integrationTestFixture.HttpClient = httpClient
				integrationTestFixture.JaegerTracer = jaegerTracer
				integrationTestFixture.Echo = echo
				integrationTestFixture.GrpcClient = grpcClient
				integrationTestFixture.ProductRepository = productRepository
				integrationTestFixture.Context = ctx
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

	//https://github.com/uber-go/fx/blob/master/app_test.go
	defer app.RequireStart().RequireStop()
	require.NoError(t, app.Err())

	configurations.ConfigMiddlewares(integrationTestFixture.Echo, integrationTestFixture.Cfg.Jaeger)

	return integrationTestFixture
}
