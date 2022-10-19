package shared

import (
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type InfrastructureConfiguration struct {
	Log               logger.ILogger
	Cfg               *config.Config
	Validator         *validator.Validate
	RabbitmqPublisher rabbitmq.IPublisher
	ConnRabbitmq      *amqp.Connection
	HttpClient        *resty.Client
	JaegerTracer      trace.Tracer
	Gorm              *gorm.DB
	Echo              *echo.Echo
}

type IdentityEndpointBase[T any] struct {
	Configuration T
	ProductsGroup *echo.Group
}
