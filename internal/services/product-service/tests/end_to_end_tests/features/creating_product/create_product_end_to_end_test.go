package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	v1_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/v1/dtos"
	test_fixture "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"net/http"
	"net/http/httptest"
	"testing"
)

type createProductEndToEndTests struct {
	*test_fixture.TestFixture
}

func TestCreateProductEndToEndTest(t *testing.T) {
	suite.Run(t, &createProductEndToEndTests{TestFixture: test_fixture.NewIntegrationTestFixture(t, fx.Options())})
}

func (c *createProductEndToEndTests) Test_Should_Return_Ok_Status_When_Create_New_Product_To_DB() {

	tsrv := httptest.NewServer(c.Echo)
	defer tsrv.Close()

	e := httpexpect.Default(c.T(), tsrv.URL)

	request := &v1_dtos.CreateProductRequestDto{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
	}

	e.POST("/api/v1/products").
		WithContext(c.Context).
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
