package test_fixture

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http"
	echserver "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	httpclient "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	gormcontainer "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/gorm_container"
	rabbitmqcontainer "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/rabbitmq_container"
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
	Ctx               context.Context
	PostgresContainer *gormcontainer.PostgresContainer
	RabbitmqContainer *rabbitmqcontainer.RabbitmqContainer
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
				http.NewContext,
				echserver.NewEchoServer,
				gormcontainer.Start,
				grpc.NewGrpcClient,
				otel.TracerProvider,
				httpclient.NewHttpClient,
				repositories.NewPostgresProductRepository,
				rabbitmqcontainer.Start,
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
				connRabbitmq *amqp.Connection,
				gormDB *gorm.DB,
				postgresContainer *gormcontainer.PostgresContainer,
				rabbitContainer *rabbitmqcontainer.RabbitmqContainer,
			) {

				integrationTestFixture.Gorm = gormDB
				integrationTestFixture.ConnRabbitmq = connRabbitmq

				integrationTestFixture.PostgresContainer = postgresContainer
				integrationTestFixture.RabbitmqContainer = rabbitContainer

				integrationTestFixture.Log = log
				integrationTestFixture.Cfg = cfg
				integrationTestFixture.RabbitmqPublisher = rabbitmqPublisher
				integrationTestFixture.HttpClient = httpClient
				integrationTestFixture.JaegerTracer = jaegerTracer
				integrationTestFixture.Echo = echo
				integrationTestFixture.GrpcClient = grpcClient
				integrationTestFixture.ProductRepository = productRepository
				integrationTestFixture.Ctx = ctx
			}),
			fx.Invoke(func(gorm *gorm.DB) error {
				return gormpgsql.Migrate(gorm, &models.Product{})
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
