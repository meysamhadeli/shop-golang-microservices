package configurations

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func InitialInfrastructures(echo *echo.Echo, log logger.ILogger, ctx context.Context, grpcClient grpc.GrpcClient,
	productRepository data.ProductRepository, config *config.Config, rabbitmqPublisher *rabbitmq.Publisher, conn *amqp.Connection,
	gorm *gorm.DB, tracer trace.Tracer, httpClient *resty.Client) *contracts.InfrastructureConfiguration {

	infar := &contracts.InfrastructureConfiguration{
		Log:               log,
		Context:           ctx,
		Echo:              echo,
		GrpcClient:        grpcClient,
		ProductRepository: productRepository,
		Cfg:               config,
		RabbitmqPublisher: rabbitmqPublisher,
		ConnRabbitmq:      conn,
		Gorm:              gorm,
		JaegerTracer:      tracer,
		HttpClient:        httpClient,
		Validator:         validator.New(),
	}

	return infar
}
