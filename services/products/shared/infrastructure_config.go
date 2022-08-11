package shared

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/interceptors"
	kafkaClient "github.com/meysamhadeli/shop-golang-microservices/pkg/kafka"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	v7 "github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type InfrastructureConfiguration struct {
	Log           logger.ILogger
	Cfg           *config.Config
	Validator     *validator.Validate
	KafkaConn     *kafka.Conn
	KafkaProducer kafkaClient.Producer
	Im            interceptors.InterceptorManager
	Gorm          *gorm.DB
	Echo          *echo.Echo
	GrpcServer    *grpc.Server
	ElasticClient *v7.Client
}

type ProductEndpointBase[T any] struct {
	Configuration T
	ProductsGroup *echo.Group
}
