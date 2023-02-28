package creating_product

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/features/creating_product/v1/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/shared/test_fixture"
	"go.uber.org/fx"
	"net/http"
	"net/http/httptest"
	"testing"
)

type createProductEndToEndTests struct {
	*test_fixture.IntegrationTestFixture
}

func TestRunner(t *testing.T) {

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("A=create-product-end-to-end-tests", func(t *testing.T) {

		var endToEndTestFixture = test_fixture.NewIntegrationTestFixture(t, fx.Options())

		testFixture := &createProductEndToEndTests{endToEndTestFixture}
		testFixture.Test_Should_Return_Ok_Status_When_Create_New_Product_To_DB()


		// Clean up the container after the test is complete
		t.Cleanup(func() {
			testFixture.PostgresContainer.Terminate(testFixture.Ctx)
			testFixture.RabbitmqContainer.Terminate(testFixture.Ctx)
		})
}

func (c *createProductEndToEndTests) Test_Should_Return_Ok_Status_When_Create_New_Product_To_DB() {

	tsrv := httptest.NewServer(c.Echo)
	defer tsrv.Close()

	e := httpexpect.Default(c.T, tsrv.URL)

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
