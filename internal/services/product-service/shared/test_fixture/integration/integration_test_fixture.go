package integration

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	gorm_test "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/gorm"
	rabbitmq_test "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/test/container/testcontainer/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

type IntegrationTestFixture struct {
	suite.Suite
	Log                logger.ILogger
	Cfg                *config.Config
	ProductRepository  data.ProductRepository
	IdentityGrpcClient grpc.GrpcClient
	Ctx                context.Context
	JaegerTracer       trace.Tracer
	RabbitMqPublisher  rabbitmq.IPublisher
	RabbitMqConsumer   rabbitmq.IConsumer
	RabbitMqConn       *amqp.Connection
}

func NewIntegrationTestFixture(t *testing.T) *IntegrationTestFixture {

	cfg, _ := config.InitConfig(constants.Test)

	log := logger.NewAppLogger(cfg.Logger)
	log.InitLogger()

	ctx, cancel := context.WithCancel(context.Background())

	err := mappings.ConfigureMappings()
	if err != nil {
		require.FailNow(t, err.Error())
	}

	cleanups := []func(){}

	identityGrpcClient, err := grpc.NewGrpcClient(cfg.IdentityGrpcClient)
	if err != nil {
		require.FailNow(t, err.Error())
	}
	cleanups = append(cleanups, func() {
		_ = identityGrpcClient.Close()
	})

	gorm, err := gorm_test.NewGormTestContainers().Start(ctx, t)
	if err != nil {
		require.FailNow(t, err.Error())
	}

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		require.FailNow(t, err.Error())
	}

	ProductRepository := repositories.NewPostgresProductRepository(log, cfg, gorm)

	tp, err := open_telemetry.TracerProvider(ctx, cfg.Jaeger, log)
	if err != nil {
		require.FailNow(t, err.Error())
	}

	jaegerTracer := tp.Tracer(cfg.Jaeger.TracerName)

	conn, err, rabbitMqCleanup := rabbitmq_test.NewRabbitMQTestContainers().Start(ctx, t)
	if err != nil {
		require.FailNow(t, err.Error())
	}

	cleanups = append(cleanups, rabbitMqCleanup)

	rabbitmqPublisher := rabbitmq.NewPublisher(cfg.Rabbitmq, conn, log, jaegerTracer)

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
		Log:                log,
		Cfg:                cfg,
		Ctx:                ctx,
		ProductRepository:  ProductRepository,
		IdentityGrpcClient: identityGrpcClient,
		JaegerTracer:       jaegerTracer,
		RabbitMqPublisher:  rabbitmqPublisher,
		RabbitMqConn:       conn,
	}

	return integration
}
