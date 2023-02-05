package main

import (
	"github.com/go-playground/validator"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http"
	echoserver "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	httpclient "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/server"
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
				gormpgsql.NewGorm,
				otel.TracerProvider,
				httpclient.NewHttpClient,
				repositories.NewPostgresInventoryRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(configurations.ConfigSwagger),
			fx.Invoke(func(gorm *gorm.DB) error {

				err := gormpgsql.Migrate(gorm, &models.Inventory{}, &models.ProductItem{})
				if err != nil {
					return err
				}
				return data.SeedData(gorm)
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigProductsMediator),
			fx.Invoke(configurations.ConfigConsumers),
		),
	).Run()
}
