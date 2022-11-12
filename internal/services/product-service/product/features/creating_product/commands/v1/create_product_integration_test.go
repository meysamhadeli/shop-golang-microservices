package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	events_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture/integration"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(t, &createProductIntegrationTests{IntegrationTestFixture: integration.NewIntegrationTestFixture(t)})
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*CreateProduct, *v1.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductId)

	createdProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, result.ProductId)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTests) BeforeTest(suiteName, testName string) {
	// some functionality before run tests
}

func (c *createProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")

	err := mediatr.RegisterRequestHandler[*CreateProduct, *v1.CreateProductResponseDto](NewCreateProductHandler(c.Log, c.Cfg, c.ProductRepository, c.RabbitMqPublisher, c.IdentityGrpcClient))
	c.Require().NoError(err)

	createProductConsumer := rabbitmq.NewConsumer(c.Cfg.Rabbitmq, c.RabbitMqConn, c.Log, c.Tracer, consumers2.HandleConsumeCreateProduct)
	err = createProductConsumer.ConsumeMessage(c.Ctx, events_v1.ProductCreated{})
	if err != nil {
		require.FailNow(c.T(), err.Error())
	}
}

func (c *createProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
