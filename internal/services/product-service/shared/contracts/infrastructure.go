package contracts

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config_options"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type InfrastructureConfiguration struct {
	Log                 logger.ILogger
	Cfg                 *config_options.Config
	Validator           *validator.Validate
	RabbitmqPublisher   rabbitmq.IPublisher
	ConnRabbitmq        *amqp.Connection
	HttpClient          *resty.Client
	JaegerTracer        trace.Tracer
	Gorm                *gorm.DB
	Echo                *echo.Echo
	GrpcClient          grpc.GrpcClient
	ProductRepository   data.ProductRepository
	Context             context.Context
	ProductEndpointBase *ProductEndpointBase[InfrastructureConfiguration]
}

type ProductEndpointBase[T any] struct {
	Configuration T
	ProductsGroup *echo.Group
}
