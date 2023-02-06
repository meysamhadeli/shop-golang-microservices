package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	creatingproductv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/commands"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/test_fixture"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/tests/unit_tests/test_data"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type createProductHandlerUnitTests struct {
	*test_fixture.UnitTestFixture
	createProductHandler *creatingproductv1commands.CreateProductHandler
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{UnitTestFixture: test_fixture.NewUnitTestFixture(t)})
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.createProductHandler = creatingproductv1commands.NewCreateProductHandler(c.Log, c.RabbitmqPublisher, c.ProductRepository, c.Ctx, c.GrpcClient)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Create_New_Product_With_Valid_Data() {

	createProductCommand := &creatingproductv1commands.CreateProduct{
		ProductID:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
		InventoryId: gofakeit.Int64(),
		Count:       gofakeit.Int32(),
	}

	product := test_data.Products[0]

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	dto, err := c.createProductHandler.Handle(c.Ctx, createProductCommand)

	c.Require().NoError(err)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.RabbitmqPublisher.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.Equal(dto.ProductId, createProductCommand.ProductID)
}
