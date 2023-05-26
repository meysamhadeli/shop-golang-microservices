package creating_product

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/consumers"
	creatingproductcommandsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/commands"
	creatingproductdtosv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/dtos"
	creatingproducteventsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/delivery"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/test_fixture"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"testing"
)

type createProductIntegrationTests struct {
	*test_fixture.IntegrationTestFixture
}

var consumer *rabbitmq.Consumer[*delivery.ProductDeliveryBase]

func TestRunner(t *testing.T) {

	var integrationTestFixture = test_fixture.NewIntegrationTestFixture(t, fx.Options(
		fx.Invoke(func(ctx context.Context, jaegerTracer trace.Tracer, log logger.ILogger, connRabbitmq *amqp.Connection, cfg *config.Config) {
			consumer = rabbitmq.NewConsumer(cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, consumers.HandleConsumeCreateProduct)
			err := consumer.ConsumeMessage(ctx, creatingproducteventsv1.ProductCreated{}, nil)
			if err != nil {
				assert.Error(t, err)
			}
		})))

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("A=create-product-integration-tests", func(t *testing.T) {

		testFixture := &createProductIntegrationTests{integrationTestFixture}
		testFixture.Test_Should_Create_New_Product_To_DB()
	})

	integrationTestFixture.PostgresContainer.Terminate(integrationTestFixture.Ctx)
	integrationTestFixture.RabbitmqContainer.Terminate(integrationTestFixture.Ctx)
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {

	command := creatingproductcommandsv1.NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000), 1, 1)
	result, err := mediatr.Send[*creatingproductcommandsv1.CreateProduct, *creatingproductdtosv1.CreateProductResponseDto](c.Ctx, command)

	assert.NoError(c.T, err)
	assert.NotNil(c.T, result)
	assert.Equal(c.T, command.ProductID, result.ProductId)

	isPublished := c.RabbitmqPublisher.IsPublished(creatingproducteventsv1.ProductCreated{})
	assert.Equal(c.T, true, isPublished)

	isConsumed := consumer.IsConsumed(creatingproducteventsv1.ProductCreated{})
	assert.Equal(c.T, true, isConsumed)

	createdProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, result.ProductId)
	assert.NoError(c.T, err)
	assert.NotNil(c.T, createdProduct)
}
