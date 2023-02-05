package test_fixture

import (
	"context"
	mocks3 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc/mocks"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	mocks2 "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq/mocks"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/constants"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product_service/tests/unit_tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"os"
	"testing"
)

type UnitTestFixture struct {
	suite.Suite
	Log               logger.ILogger
	Cfg               *config.Config
	Ctx               context.Context
	RabbitmqPublisher *mocks2.IPublisher
	RabbitmqConsumer  *mocks2.IConsumer
	ProductRepository *mocks.ProductRepository
	GrpcClient        *mocks3.GrpcClient
}

func NewUnitTestFixture(t *testing.T) *UnitTestFixture {

	err := os.Setenv("APP_ENV", constants.Test)

	if err != nil {
		require.FailNow(t, err.Error())
	}

	unitTestFixture := &UnitTestFixture{}

	app := fxtest.New(t,
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				http.NewContext,
			),
			fx.Invoke(func(
				ctx context.Context,
				log logger.ILogger,
				cfg *config.Config,
			) {
				unitTestFixture.Log = log
				unitTestFixture.Cfg = cfg
				unitTestFixture.Ctx = ctx
			}),
			fx.Invoke(mappings.ConfigureMappings),
		),
	)

	//https://github.com/uber-go/fx/blob/master/app_test.go
	defer app.RequireStart().RequireStop()
	require.NoError(t, app.Err())

	// create new mocks
	unitTestFixture.RabbitmqPublisher = &mocks2.IPublisher{}
	unitTestFixture.RabbitmqConsumer = &mocks2.IConsumer{}
	unitTestFixture.ProductRepository = &mocks.ProductRepository{}
	unitTestFixture.GrpcClient = &mocks3.GrpcClient{}

	unitTestFixture.RabbitmqPublisher.On("PublishMessage", mock.Anything, mock.Anything).Return(nil)
	unitTestFixture.RabbitmqConsumer.On("ConsumeMessage", mock.Anything, mock.Anything).Return(nil)

	return unitTestFixture
}
