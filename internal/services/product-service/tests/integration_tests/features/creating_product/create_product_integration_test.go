package creating_product

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/commands"
	v1_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/dtos"
	v1_event "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"testing"
)

type createProductIntegrationTests struct {
	*test_fixture.TestFixture
}

var consumer *rabbitmq.Consumer

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(t, &createProductIntegrationTests{TestFixture: test_fixture.NewIntegrationTestFixture(t, fx.Options(
		fx.Invoke(func(ctx context.Context, jaegerTracer trace.Tracer, log logger.ILogger, connRabbitmq *amqp.Connection, cfg *config.Config) {
			consumer = rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers2.HandleConsumeCreateProduct)
			err := consumer.ConsumeMessage(ctx, v1_event.ProductCreated{})
			if err != nil {
				require.FailNow(t, err.Error())
			}
		}),
	))})
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {

	command := commands.NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*commands.CreateProduct, *v1_dtos.CreateProductResponseDto](c.Context, command)
	c.Require().NoError(err)

	isPublished := c.RabbitmqPublisher.IsPublished(v1_event.ProductCreated{})
	c.Assert().Equal(true, isPublished)

	isConsumed := consumer.IsConsumed(v1_event.ProductCreated{})
	c.Assert().Equal(true, isConsumed)

	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductId)

	createdProduct, err := c.TestFixture.ProductRepository.GetProductById(c.Context, result.ProductId)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTests) BeforeTest(suiteName, testName string) {
	// some functionality before run tests
}

func (c *createProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
}

func (c *createProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
