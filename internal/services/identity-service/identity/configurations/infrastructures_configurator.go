package configurations

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	contracts2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func InitialInfrastructures(echo *echo.Echo, log logger.ILogger, ctx context.Context, grpcServer *grpc.GrpcServer,
	userRepository contracts2.UserRepository, config *config.Config, rabbitmqPublisher *rabbitmq.Publisher, conn *amqp.Connection,
	gorm *gorm.DB, tracer trace.Tracer, httpClient *resty.Client) *contracts.InfrastructureConfiguration {

	infar := &contracts.InfrastructureConfiguration{
		Log:               log,
		Context:           ctx,
		Echo:              echo,
		GrpcServer:        grpcServer,
		UserRepository:    userRepository,
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
