package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	v1_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture/integration"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type createProductEndToEndTests struct {
	*integration.IntegrationTestFixture
}

func TestCreateProductEndToEndTest(t *testing.T) {
	suite.Run(t, &createProductEndToEndTests{IntegrationTestFixture: integration.NewIntegrationTestFixture(t)})
}

func (c *createProductEndToEndTests) Test_Should_Create_New_Product_To_DB() {

	//createProductEndpoint := creating_product.NewCreteProductEndpoint(c.ProductEndpointBase)
	//createProductEndpoint.MapRoute()

	//c.Run()

	request := &v1_dtos.CreateProductRequestDto{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
	}

	// create httpexpect instance
	expect := httpexpect.Default(c.T(), c.Cfg.Echo.BasePathAddress())

	expect.POST("").
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
