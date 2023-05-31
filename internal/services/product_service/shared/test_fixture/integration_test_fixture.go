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
	gormcontainer "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/postgres_container"
	rabbitmqcontainer "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/rabbitmq_container"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/data/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/delivery"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
	"os"
	"testing"
)

type IntegrationTestFixture struct {
	*testing.T
	Log               logger.ILogger
	Cfg               *config.Config
	RabbitmqPublisher rabbitmq.IPublisher
	RabbitmqConsumer  *rabbitmq.Consumer[delivery.ProductDeliveryBase]
	ConnRabbitmq      *amqp.Connection
	HttpClient        *resty.Client
	JaegerTracer      trace.Tracer
	Gorm              *gorm.DB
	Echo              *echo.Echo
	GrpcClient        grpc.GrpcClient
	ProductRepository contracts.ProductRepository
	Ctx               context.Context
}

func NewIntegrationTestFixture(t *testing.T, option fx.Option) *IntegrationTestFixture {

	err := os.Setenv("APP_ENV", constants.Test)

	if err != nil {
		require.FailNow(t, err.Error())
	}

	integrationTestFixture := &IntegrationTestFixture{T: t}

	app := fxtest.New(t,
		fx.Options(
			fx.Provide(
				func() *testing.T {
					return t
				},
				http.NewContext,
				config.InitConfig,
				rabbitmqcontainer.Start,
				gormcontainer.Start,
				logger.InitLogger,
				echserver.NewEchoServer,
				grpc.NewGrpcClient,
				otel.TracerProvider,
				httpclient.NewHttpClient,
				repositories.NewPostgresProductRepository,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(func(
				rabbitmqPublisher rabbitmq.IPublisher,
				productRepository contracts.ProductRepository,
				grpcClient grpc.GrpcClient,
				echo *echo.Echo,
				log logger.ILogger,
				jaegerTracer trace.Tracer,
				httpClient *resty.Client,
				validator *validator.Validate,
				cfg *config.Config,
				connRabbitmq *amqp.Connection,
				gormDB *gorm.DB,
				ctx context.Context,
			) {
				integrationTestFixture.Gorm = gormDB
				integrationTestFixture.ConnRabbitmq = connRabbitmq

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

	// Start the Uber FX application
	if err := app.Start(integrationTestFixture.Ctx); err != nil {
		t.Fatalf("failed to start the Uber FX application: %v", err)
	}

	defer func(app *fxtest.App, ctx context.Context) {
		err := app.Stop(ctx)
		if err != nil {
			t.Fatalf("failed to stop the Uber FX application: %v", err)
		}
	}(app, integrationTestFixture.Ctx)

	configurations.ConfigMiddlewares(integrationTestFixture.Echo, integrationTestFixture.Cfg.Jaeger)

	return integrationTestFixture
}
