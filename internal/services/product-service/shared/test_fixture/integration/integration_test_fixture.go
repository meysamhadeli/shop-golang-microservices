package integration

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

type IntegrationTestFixture struct {
	suite.Suite
	Log                 logger.ILogger
	Cfg                 *config.Config
	ProductRepository   data.ProductRepository
	IdentityGrpcClient  grpc.GrpcClient
	Ctx                 context.Context
	JaegerTracer        trace.Tracer
	RabbitMqPublisher   rabbitmq.IPublisher
	RabbitMqConsumer    rabbitmq.IConsumer
	RabbitMqConn        *amqp.Connection
	EchoServer          *server.EchoServer
	ProductEndpointBase *contracts.ProductEndpointBase[contracts.InfrastructureConfiguration]
}

func NewIntegrationTestFixture(t *testing.T) *IntegrationTestFixture {

	c := NewTestInfrastructureConfigurator(t)

	infrastructures, echoServer, err, infraCleanup := c.ConfigInfrastructures()

	if err != nil {
		require.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		// with Cancel() we send signal to done() channel to stop  grpc, http and workers gracefully
		//https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		//https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		mediatr.ClearRequestRegistrations()
		infraCleanup()
	})

	integration := &IntegrationTestFixture{
		Log:                 infrastructures.Log,
		Cfg:                 infrastructures.Cfg,
		Ctx:                 infrastructures.Context,
		ProductRepository:   infrastructures.ProductRepository,
		IdentityGrpcClient:  infrastructures.GrpcClient,
		JaegerTracer:        infrastructures.JaegerTracer,
		RabbitMqPublisher:   infrastructures.RabbitmqPublisher,
		RabbitMqConn:        infrastructures.ConnRabbitmq,
		ProductEndpointBase: infrastructures.ProductEndpointBase,
		EchoServer:          echoServer,
	}

	return integration
}
