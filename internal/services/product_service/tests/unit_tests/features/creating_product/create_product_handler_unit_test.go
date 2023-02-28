package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	creatingproductv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/commands"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/test_fixture"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/tests/unit_tests/test_data"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type createProductHandlerUnitTests struct {
	*test_fixture.UnitTestFixture
	createProductHandler *creatingproductv1commands.CreateProductHandler
}

func TestRunner(t *testing.T) {

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("A=create-product-unit-tests", func(t *testing.T) {

		var unitTestFixture = test_fixture.NewUnitTestFixture(t)

		mockCreateProductHandler := creatingproductv1commands.NewCreateProductHandler(unitTestFixture.Log, unitTestFixture.RabbitmqPublisher,
			unitTestFixture.ProductRepository, unitTestFixture.Ctx, unitTestFixture.GrpcClient)

		testFixture := &createProductHandlerUnitTests{unitTestFixture, mockCreateProductHandler}
		testFixture.Test_Handle_Should_Create_New_Product_With_Valid_Data()
	})
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
		InventoryId: 1,
		Count:       1,
	}

	product := test_data.Products[0]

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	dto, err := c.createProductHandler.Handle(c.Ctx, createProductCommand)

	assert.NoError(c.T, err)

	c.ProductRepository.AssertNumberOfCalls(c.T, "CreateProduct", 1)
	c.RabbitmqPublisher.AssertNumberOfCalls(c.T, "PublishMessage", 1)
	assert.Equal(c.T, dto.ProductId, createProductCommand.ProductID)
}
