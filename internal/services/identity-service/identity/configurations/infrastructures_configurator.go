package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
	"google.golang.org/grpc"
	"net/http"
)

type infrastructureConfigurator struct {
	Log        logger.ILogger
	Cfg        *config.Config
	Echo       *echo.Echo
	GrpcServer *grpc.Server
}

func NewInfrastructureConfigurator(log logger.ILogger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *infrastructureConfigurator {
	return &infrastructureConfigurator{Cfg: cfg, Echo: echo, Log: log, GrpcServer: grpcServer}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, func()) {

	infrastructure := &contracts.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, Log: ic.Log, Validator: validator.New()}

	cleanups := []func(){}

	gorm, err := gorm_postgres.NewGorm(ic.Cfg.GormPostgres)

	if err != nil {
		return err, nil
	}
	infrastructure.Gorm = gorm

	err = gorm.AutoMigrate(&models.User{})
	if err != nil {
		return err, nil
	}

	tp, err := open_telemetry.TracerProvider(ctx, ic.Cfg.Jaeger, ic.Log)
	if err != nil {
		ic.Log.Fatal(err)
		return err, nil
	}

	infrastructure.JaegerTracer = tp.Tracer(ic.Cfg.Jaeger.TracerName)

	ic.Log.Infof("%s is running", config.GetMicroserviceName(ic.Cfg.ServiceName))

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	configSwagger(ic.Echo)

	configMiddlewares(ic.Echo, ic.Cfg.Jaeger)

	configureOauth2(ic.Echo)

	ConfigIdentityGrpcServer(ctx, ic.GrpcServer, infrastructure)

	pc := NewUsersModuleConfigurator(infrastructure)

	err = pc.ConfigureIdentitiesModule(ctx)
	if err != nil {
		return err, nil
	}

	ic.Echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", config.GetMicroserviceName(ic.Cfg.ServiceName)))
	})

	return nil, func() {
		for _, c := range cleanups {
			c()
		}
	}
}
