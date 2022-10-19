package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type infrastructureConfigurator struct {
	Log  logger.ILogger
	Cfg  *config.Config
	Echo *echo.Echo
}

func NewInfrastructureConfigurator(log logger.ILogger, cfg *config.Config, echo *echo.Echo) *infrastructureConfigurator {
	return &infrastructureConfigurator{Cfg: cfg, Echo: echo, Log: log}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, func()) {

	infrastructure := &shared.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, Log: ic.Log, Validator: validator.New()}

	cleanups := []func(){}

	gorm, err := ic.configGorm()
	if err != nil {
		return err, nil
	}
	infrastructure.Gorm = gorm

	tp, err := ic.configOpenTelemetry()
	if err != nil {
		return err, nil
	}
	infrastructure.JaegerTracer = tp.Tracer(ic.Cfg.Jaeger.TracerName)

	cleanups = append(cleanups, func() {
		err = tp.Shutdown(ctx)
		if err != nil {
			ic.Log.Fatal(err)
		}
	})

	ic.Log.Infof("%s is running", config.GetMicroserviceName(ic.Cfg.ServiceName))

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	ic.configSwagger()

	ic.configMiddlewares(ic.Cfg.Jaeger)
	if err != nil {
		return err, nil
	}

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
