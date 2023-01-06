package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	v1_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture/integration"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"net/http/httptest"
	"testing"
)

type createProductEndToEndTests struct {
	*integration.IntegrationTestFixture
}

func TestCreateProductEndToEndTest(t *testing.T) {
	suite.Run(t, &createProductEndToEndTests{IntegrationTestFixture: integration.NewIntegrationTestFixture(t, fx.Options())})
}

func (c *createProductEndToEndTests) Test_Should_Return_Ok_Status_When_Create_New_Product_To_DB() {

	s := httptest.NewServer(c.Echo)
	defer s.Close()

	request := &v1_dtos.CreateProductRequestDto{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
	}

	// create httpexpect instance
	expect := httpexpect.Default(c.T(), s.URL)

	expect.POST("api/v1/products").
		WithContext(c.Context).
		WithJSON(request)
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
