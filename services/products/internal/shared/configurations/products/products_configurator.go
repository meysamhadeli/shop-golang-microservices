package products

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/configurations"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/configurations/infrastructure"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/web"
	"google.golang.org/grpc"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type catalogsServiceConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echo       *echo.Echo
	grpcServer *grpc.Server
}

func NewCatalogsServiceConfigurator(log logger.Logger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *catalogsServiceConfigurator {
	return &catalogsServiceConfigurator{cfg: cfg, echo: echo, grpcServer: grpcServer, log: log}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (error, func()) {

	ic := infrastructure.NewInfrastructureConfigurator(c.log, c.cfg, c.echo, c.grpcServer)
	infrastructureConfigurations, err, infraCleanup := ic.ConfigInfrastructures(ctx)
	if err != nil {
		return err, nil
	}

	//------------------------------------------------------------------------------//

	pc := configurations.NewProductsModuleConfigurator(infrastructureConfigurations)
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil
	}

	err = c.migrateCatalogs(infrastructureConfigurations.Gorm)
	if err != nil {
		return err, nil
	}

	c.echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.cfg)))
	})

	return nil, infraCleanup
}
