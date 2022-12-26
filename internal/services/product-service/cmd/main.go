package main

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/context_provider"
	echo_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	fx.New(
		fx.Provide(
			config.InitConfig,
			logger.InitLogger,
			context_provider.NewContext,
			echo_server.NewEchoServer,
			grpc.NewGrpcClient,
			gorm_postgres.NewGorm,
			open_telemetry.TracerProvider,
			http_client.NewHttpClient,
			repositories.NewPostgresProductRepository,
			rabbitmq.NewRabbitMQConn,
			rabbitmq.NewPublisher,
			configurations.InitialInfrastructures,
		),
		fx.Invoke(echo_server.RunEchoServer),
		fx.Invoke(configurations.ConfigMiddlewares),
		fx.Invoke(configurations.ConfigSwagger),
		fx.Invoke(func(gorm *gorm.DB) error {
			return gorm_postgres.Migrate(gorm, &models.Product{})
		}),
		fx.Invoke(mappings.ConfigureMappings),
		fx.Invoke(configurations.ConfigEndpoints),
		fx.Invoke(configurations.ConfigProductsMediator),
		fx.Invoke(configurations.ConfigConsumers),
	).Run()
}
