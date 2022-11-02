package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	repositories_imp "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*shared.InfrastructureConfiguration
	IdentityGrpcClient grpc.GrpcClient
}

func NewProductsModuleConfigurator(infrastructure *shared.InfrastructureConfiguration, identityGrpcClient grpc.GrpcClient) *productsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfiguration: infrastructure, IdentityGrpcClient: identityGrpcClient}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/products")

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

	return nil
}
