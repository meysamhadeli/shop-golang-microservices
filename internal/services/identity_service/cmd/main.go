package main

import (
	"github.com/go-playground/validator"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http"
	echoserver "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/oauth2"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/server"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				http.NewContext,
				echoserver.NewEchoServer,
				grpc.NewGrpcServer,
				gormpgsql.NewGorm,
				otel.TracerProvider,
				httpclient.NewHttpClient,
				repositories.NewPostgresUserRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(configurations.ConfigSwagger),
			fx.Invoke(func(gorm *gorm.DB) error {
				return gormpgsql.Migrate(gorm, &models.User{})
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigUsersMediator),
			fx.Invoke(oauth2.RunOauthServer),
		),
	).Run()
}
