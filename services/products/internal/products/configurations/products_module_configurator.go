package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	repositories_imp "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/mappings"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*config.InfrastructureConfiguration
}

func NewProductsModuleConfigurator(infrastructure *config.InfrastructureConfiguration) *productsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfiguration: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/" + c.Cfg.Http.ProductsPath)

	productRepository := repositories_imp.NewPostgresProductRepository(c.Log, c.Cfg, c.Gorm)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = c.configProductsMediator(productRepository)
	if err != nil {
		return err
	}

	c.configEndpoints(ctx, group)

	//if c.Cfg.DeliveryType == "grpc" {
	//	c.configGrpc(ctx)
	//}

	return nil
}
