package main

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/context_provider"
	echo_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/oauth2"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	configurations "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/configurations"
	repositories "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/server"
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
			grpc.NewGrpcServer,
			gorm_postgres.NewGorm,
			open_telemetry.TracerProvider,
			http_client.NewHttpClient,
			repositories.NewPostgresUserRepository,
			rabbitmq.NewRabbitMQConn,
			rabbitmq.NewPublisher,
			configurations.InitialInfrastructures,
		),
		fx.Invoke(server.RunServers),
		fx.Invoke(configurations.ConfigMiddlewares),
		fx.Invoke(configurations.ConfigSwagger),
		fx.Invoke(func(gorm *gorm.DB) error {
			return gorm_postgres.Migrate(gorm, &models.User{})
		}),
		fx.Invoke(mappings.ConfigureMappings),
		fx.Invoke(configurations.ConfigEndpoints),
		fx.Invoke(configurations.ConfigUsersMediator),
		fx.Invoke(oauth2.RunOauthServer),
	).Run()
}
