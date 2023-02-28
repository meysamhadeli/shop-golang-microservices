package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/test_fixture"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"net/http"
	"net/http/httptest"
	"testing"
)

type createProductEndToEndTests struct {
	*test_fixture.IntegrationTestFixture
}

func TestCreateProductEndToEndTest(t *testing.T) {
	suite.Run(t, &createProductEndToEndTests{IntegrationTestFixture: test_fixture.NewIntegrationTestFixture(t, fx.Options())})
}

func (c *createProductEndToEndTests) Test_Should_Return_Ok_Status_When_Create_New_Product_To_DB() {

	defer c.PostgresContainer.Terminate(c.Ctx)
	defer c.RabbitmqContainer.Terminate(c.Ctx)

	tsrv := httptest.NewServer(c.Echo)
	defer tsrv.Close()

	e := httpexpect.Default(c.T(), tsrv.URL)

	request := &dtos.CreateProductRequestDto{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
		InventoryId: 1,
		Count:       1,
	}

	e.POST("/api/v1/products").
		WithContext(c.Ctx).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}

func (c *createProductEndToEndTests) BeforeTest(suiteName, testName string) {
	// some functionality before run tests
}

func (c *createProductEndToEndTests) SetupTest() {
	c.T().Log("SetupTest")
}

func (c *createProductEndToEndTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
