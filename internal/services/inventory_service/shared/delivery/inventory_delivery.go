package delivery

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data/contracts"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type InventoryDeliveryBase struct {
	Log                 logger.ILogger
	Cfg                 *config.Config
	RabbitmqPublisher   rabbitmq.IPublisher
	ConnRabbitmq        *amqp.Connection
	HttpClient          *resty.Client
	JaegerTracer        trace.Tracer
	Gorm                *gorm.DB
	Echo                *echo.Echo
	GrpcClient          grpc.GrpcClient
	InventoryRepository contracts.InventoryRepository
	Ctx                 context.Context
}
